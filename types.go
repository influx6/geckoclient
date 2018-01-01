package geckoclient

// DataType defines a interface type which exposes a single method
// for returning field values of type details. It helps to construct
// the values sent along a NewDataset field list for items to be created
// on the geckoboard API.
type DataType interface {
	Field() map[string]interface{}
}

// AnyType defines a type alias for a map, which can be used to
// wrap said map to match the DataType interface.
// It implements the DataType interface on the map type.
type AnyType map[string]interface{}

// Field implements the DataType interface.
func (a AnyType) Field() map[string]interface{} {
	return a
}

// DateTimeType representing the supported datetime type supported by the geckoboard API.
// The provided name and optional value are used to generate appropriate map
// of values expected by the API. Values to be used in data must be formatted in ISO 8601.
type DateTimeType struct {
	Name string
}

// Field implements the DataType interface.
func (d DateTimeType) Field() map[string]interface{} {
	return map[string]interface{}{
		"name": d.Name,
		"type": "datetime",
	}
}

// DateType representing the supported date type supported by the geckoboard API.
// The provided name and optional value are used to generate appropriate map
// of values expected by the API.
type DateType struct {
	Name string
}

// Field implements the DataType interface.
func (d DateType) Field() map[string]interface{} {
	return map[string]interface{}{
		"name": d.Name,
		"type": "date",
	}
}

// StringType representing the supported string type supported by the geckoboard API.
// The provided name and optional value are used to generate appropriate map
// of values expected by the API.
type StringType struct {
	Name string
}

// Field implements the DataType interface.
func (d StringType) Field() map[string]interface{} {
	return map[string]interface{}{
		"name": d.Name,
		"type": "string",
	}
}

// NumberType representing the supported number type supported by the geckoboard API.
// The provided name and optional value are used to generate appropriate map
// of values expected by the API.
type NumberType struct {
	Name     string
	Optional bool
}

// Field implements the DataType interface.
func (d NumberType) Field() map[string]interface{} {
	return map[string]interface{}{
		"name":     d.Name,
		"type":     "number",
		"optional": d.Optional,
	}
}

// MoneyType representing the supported money type supported by the geckoboard API.
// The provided name and optional value are used to generate appropriate map
// of values expected by the API.
// Currency code must be in abbreviations where USD representings United State Dollars.
// Values to be used for this type must be in smallest denomination values with respect
// to currency, where if USD for currency code, then value of $10.00 should be written in
// 10000 cents.
type MoneyType struct {
	Name         string
	CurrencyCode string
	Optional     bool
}

// Field implements the DataType interface.
func (d MoneyType) Field() map[string]interface{} {
	return map[string]interface{}{
		"name":          d.Name,
		"type":          "money",
		"optional":      d.Optional,
		"currency_code": d.CurrencyCode,
	}
}

// PercentageType representing the supported percentage type supported by the geckoboard API.
// The provided name and optional value are used to generate appropriate map
// of values expected by the API. Values used for this type must be between 0 and 1, i.e
// 0.1 which gets turned to 10% ...etc.
type PercentageType struct {
	Name     string
	Optional bool
}

// Field implements the DataType interface.
func (d PercentageType) Field() map[string]interface{} {
	return map[string]interface{}{
		"name":     d.Name,
		"type":     "percentage",
		"optional": d.Optional,
	}
}
