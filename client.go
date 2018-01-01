package geckoclient

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

var (
	api     = "https://api.geckoboard.com"
	hclient = &http.Client{
	//Timeout: time.Second * 30,
	}
)

// errors ...
var (
	ErrInvalidRequest      = errors.New("request is invalid: api respond as bad request")
	ErrRequestConflict     = errors.New("request encountered resource conflict")
	ErrFailedRequest       = errors.New("request failed: unknown status response")
	ErrBadCredentials      = errors.New("request denied due to bad auth crendentials")
	ErrInvalidResponseType = errors.New("invalid response type, expected 'application/json'")
)

// NewDataset embodies the data sent as json to create a new
// data set within user associated account.
type NewDataset struct {
	Fields   map[string]DataType
	UniqueBy []string
}

// Dataset embodies the data send along a http request to either add, update or replace
// existing dataset data for a existing dataset within user's Geckoboard account.
type Dataset struct {
	// Data to be added to or replace existing dataset data.
	Data []map[string]interface{} `json:"data"`

	// Provide field names of dataset to use as deletion criteria.
	// Only used when adding or appending more data, not when replacing datasets
	// data.
	DeleteBy []string `json:"delete_by,omitempty"`
}

// APIError embodies the data received from the GeckoBoard API when
// a request returns an associated error response.
type APIError struct {
	Message string `json:"message"`
}

// Error returns associated error associated with the instance.
func (a APIError) Error() string {
	return a.Message
}

// Client embodies a http client for interacting with the GeckoBoard API.
type Client struct {
	auth   string
	agent  string
	apiURL string
}

// New returns a new instance of a Client for interacting with
// the Geckoboard API.
func New(authKey string) (Client, error) {
	return CustomClient(api, authKey, "")
}

// NewWithUserAgent returns a new instance of a Client for interacting with
// the Geckoboard API.
func NewWithUserAgent(authKey string, agent string) (Client, error) {
	return CustomClient(api, authKey, agent)
}

// CustomClient returns a new instance of a Client for interacting with
// the Geckoboard API found from provided apiURL provided with auth key and agent
// name if provided.
func CustomClient(apiURL string, authKey string, agent string) (Client, error) {
	var gc Client
	gc.agent = agent
	gc.auth = authKey
	gc.apiURL = apiURL
	return gc, gc.verify()
}

// ReplaceData replaces all data with provided set for giving datasetID if it exists within user's
// Geckoboard API account based on Auth key.
//
// PUT https://api.geckoboard.com/datasets/:datasetid/data
//
// Note: Dataset.DeleteBy is ignored and not used during this call.
func (gc Client) ReplaceData(ctx context.Context, datasetID string, data Dataset) error {
	var body bytes.Buffer
	if err := json.NewEncoder(&body).Encode(data); err != nil {
		return err
	}

	res, err := gc.doRequest(ctx, "PUT", fmt.Sprintf("datasets/%s/data", datasetID), &body)
	if err != nil {
		return err
	}

	defer res.Close()
	return nil
}

// PushData adds new data into the data set uploaded to the Geckoboard API represented by
// the provided ID. The data are added if new and will be merged based on standards of the
// unique_by field names if present within dataset.
//
// POST https://api.geckoboard.com/datasets/:datasetid/data
//
// NOTE: To replace use Client.ReplaceData.
// NOTE: Dataset.DeleteBy is used if provided, during this call.
func (gc Client) PushData(ctx context.Context, datasetID string, data Dataset) error {
	var body bytes.Buffer
	if err := json.NewEncoder(&body).Encode(data); err != nil {
		return err
	}

	res, err := gc.doRequest(ctx, "POST", fmt.Sprintf("datasets/%s/data", datasetID), &body)
	if err != nil {
		return err
	}

	defer res.Close()
	return nil
}

// Delete sends a request to delete dataset marked by provided datasetID. It issues a DELETE request
// to the Geckoboard API.
//
// DELETE https://api.geckoboard.com/datasets/:datasetid
//
func (gc Client) Delete(ctx context.Context, datasetID string) error {
	res, err := gc.doRequest(ctx, "PUT", fmt.Sprintf("datasets/%s", datasetID), nil)
	if err != nil {
		return err
	}

	defer res.Close()
	return nil
}

// Create sends a request to create the provided dataset by issuing a http request
// to the Geckoboard API.
//
// PUT https://api.geckoboard.com/datasets/:datasetid
//
func (gc Client) Create(ctx context.Context, datasetID string, set NewDataset) error {
	newData := struct {
		Fields   map[string]interface{} `json:"fields"`
		UniqueBy []string               `json:"unique_by,omitempty"`
	}{
		Fields: map[string]interface{}{},
	}

	newData.UniqueBy = set.UniqueBy
	for name, field := range set.Fields {
		newData.Fields[name] = field.Field()
	}

	var body bytes.Buffer
	if err := json.NewEncoder(&body).Encode(newData); err != nil {
		return err
	}

	res, err := gc.doRequest(ctx, "PUT", fmt.Sprintf("datasets/%s", datasetID), &body)
	if err != nil {
		return err
	}

	defer res.Close()
	return nil
}

// verify validates authenticity of api token for giving client.
func (gc Client) verify() error {
	_, err := gc.doRequest(context.Background(), "GET", "/", nil)
	return err
}

// doRequest contains necessary logic to send request to API endpoint and appropriately return
// desired response.
func (gc Client) doRequest(ctx context.Context, method string, path string, body io.Reader) (io.ReadCloser, error) {
	req, err := http.NewRequest(method, fmt.Sprintf("%s/%s", gc.apiURL, path), body)
	if err != nil {
		return nil, err
	}

	if ctx != nil {
		req = req.WithContext(ctx)
	}

	req.SetBasicAuth(gc.auth, "")
	req.Header.Set("Content-Type", "application/json")

	if gc.agent != "" {
		req.Header.Set("User-Agent", gc.agent)
	}

	res, err := hclient.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode >= 200 && res.StatusCode < 300 {
		return res.Body, nil
	}

	if res.StatusCode >= 400 {
		if res.StatusCode == http.StatusConflict {
			return nil, ErrRequestConflict
		}

		if res.StatusCode == http.StatusBadRequest {
			return nil, ErrInvalidRequest
		}

		if res.StatusCode == http.StatusUnauthorized {
			return nil, ErrBadCredentials
		}

		return nil, ErrFailedRequest
	}

	if !strings.Contains(res.Header.Get("Content-Type"), "application/json") {
		return nil, ErrInvalidResponseType
	}

	defer res.Body.Close()

	var recErr = struct {
		Error APIError `json:"error"`
	}{}

	if err := json.NewDecoder(res.Body).Decode(&recErr); err != nil {
		return nil, ErrFailedRequest
	}

	return nil, recErr.Error
}
