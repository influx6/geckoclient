# GeckoClient
Go/Golang client library for Geckoboard API (https://developer.geckoboard.com/api-reference/).

## Install

```bash
go get github.com/influx6/geckoclient
```

### Create new client

New Client verifies API Key to ensure user has valid API auth.

```go
client, err := geckoclient.New("222efc82e7933138077b1c2554439e15")
```

### Create Dataset

Create new dataset.

```go
client.Create(context.Background(), "sales.gross", geckoclient.NewDataset{
	UniqueBy: []string{"date"},
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
```

`unique_by` is an optional array of one or more field names whose values will be unique across all your records.

Available field types:

- `DateType`
- `DateTimeType`
- `NumberType`
- `PercentageType`
- `StringType`
- `MoneyType`


### Delete

Delete a dataset and all data therein.

```go
client.Delete(context.Background(), "sales.gross")
```

### Put

Replace all data in the dataset.

```go
client.ReplaceData(context.Background(), "sales.gross", geckoclient.Dataset{
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
```

### Post

Append data to a dataset.

```go
client.PushData(context.Background(), "sales.gross", geckoclient.Dataset{
	DeleteBy: []string{"name"},
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
```

`delete_by` is an optional field by which to order the truncation of records once the maximum record count has been reached. By default the oldest records (by insertion time) will be removed.

## Vendoring

Vendoring is done with [Govendor](https://github.com/kardianos/govendor).

## Contributing

1. Fork it ( https://github.com/influx6/geckoclient/fork )
2. Create your feature branch (`git checkout -b my-new-feature`)
3. Commit your changes (`git commit -am 'Add some feature'`)
4. Push to the branch (`git push origin my-new-feature`)
5. Create a new Pull Request