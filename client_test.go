package geckoclient_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/geckoboard/geckoclient"
	"github.com/influx6/faux/tests"
)

var (
	now     = time.Now()
	isoTime = "2016-01-01"
	apiAuth = os.Getenv("GECKOBOARD_TEST_API_KEY")
)

func TestGeckoClient(t *testing.T) {
	tests.Info("Running tests with APIKey %+q", apiAuth)
	client, err := geckoclient.New(apiAuth)
	if err != nil {
		tests.FailedWithError(err, "Should have successfully created new geckoboard api client")
	}
	tests.Passed("Should have successfully created new geckoboard api client")

	testDatasetCreate(t, client)
	testDatasetPushData(t, client)
	testDatasetReplaceData(t, client)
	testDatasetDelete(t, client)
}

func testDatasetCreate(t *testing.T, client geckoclient.Client) {
	tests.Header("When creating a new dataset")
	{
		err := client.Create(context.Background(), "decking", geckoclient.NewDataset{
			Fields: map[string]geckoclient.DataType{
				"name": geckoclient.StringType{
					Name: "transaction_target",
				},
				"amount": geckoclient.NumberType{
					Name: "transaction_amount",
				},
				"date": geckoclient.DateType{
					Name: "transaction_date",
				},
			},
		})

		if err != nil {
			tests.FailedWithError(err, "Should have successfully created new dataset")
		}
		tests.Passed("Should have successfully created new dataset")
	}
}

func testDatasetPushData(t *testing.T, client geckoclient.Client) {
	tests.Header("When pushing data in a dataset")
	{
		err := client.PushData(context.Background(), "decking", geckoclient.Dataset{
			Data: []map[string]interface{}{
				{
					"amount": 300,
					"name":   "Waxon Butter",
					"date":   now.Format(isoTime),
				},
				{
					"amount": 1300,
					"name":   "Shred Lack",
					"date":   now.Add(time.Hour * 30).Format(isoTime),
				},
				{
					"amount": 500,
					"name":   "Creg Washer",
					"date":   now.Add(time.Hour * 10).Format(isoTime),
				},
			},
		})

		if err != nil {
			tests.FailedWithError(err, "Should have successfully added new data to dataset")
		}
		tests.Passed("Should have successfully added new data to dataset")
	}
}

func testDatasetReplaceData(t *testing.T, client geckoclient.Client) {
	tests.Header("When replacing data in a dataset")
	{
		err := client.ReplaceData(context.Background(), "decking", geckoclient.Dataset{
			Data: []map[string]interface{}{
				{
					"amount": 300,
					"name":   "Waxon Rutter",
					"date":   time.Now().Format(isoTime),
				},
				{
					"amount": 1300,
					"name":   "Shred Hack",
					"date":   now.Add(time.Hour * 10).Format(isoTime),
				},
				{
					"amount": 500,
					"name":   "Creg Washer",
					"date":   now.Add(time.Hour * 50).Format(isoTime),
				},
			},
		})

		if err != nil {
			tests.FailedWithError(err, "Should have successfully replaced old data with new data in dataset")
		}
		tests.Passed("Should have successfully replaced old data with new data in dataset")
	}
}

func testDatasetDelete(t *testing.T, client geckoclient.Client) {
	tests.Header("When deleting a new dataset")
	{
		err := client.Delete(context.Background(), "decking")
		if err != nil {
			tests.FailedWithError(err, "Should have successfully deleted new dataset")
		}
		tests.Passed("Should have successfully deleted new dataset")
	}
}
