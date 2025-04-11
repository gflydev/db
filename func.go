package db

import (
	"encoding/json"
	"github.com/gflydev/core/errors"
	"reflect"
	"strconv"
)

// setValue uses reflection to set a value for a key in a struct.
// Supported types: bool, string, int (Int, Int8, Int16, Int32, Int64),
// uint (UInt, UInt8, UInt16, UInt32, UInt64), float (float32, float64)
//
// Parameters:
//   - model (any): The struct or pointer to a struct on which the value will be set.
//   - key (string): The field name in the struct to set the value for.
//   - data (any): The value to set in the specified field.
//
// Returns:
//   - error: Returns an error if the operation fails.
func setValue(model any, key string, data any) (err error) {
	// Get the reflect.Value of the model
	value := reflect.ValueOf(model)

	// Retrieve the field by name
	field := reflect.Indirect(value).FieldByName(key)

	var val reflect.Value

	// Convert the data to a string representation
	dataStr := toStr(data)

	switch field.Kind() {
	case reflect.String:
		val = reflect.ValueOf(data.(string))
	case reflect.Bool:
		val = reflect.ValueOf(data.(bool))
	case reflect.Int:
		intVar, _ := strconv.Atoi(dataStr)
		val = reflect.ValueOf(intVar)
	case reflect.Int8:
		intVar, _ := strconv.ParseInt(dataStr, 10, 8)
		val = reflect.ValueOf(int8(intVar))
	case reflect.Int16:
		intVar, _ := strconv.ParseInt(dataStr, 10, 16)
		val = reflect.ValueOf(int16(intVar))
	case reflect.Int32:
		intVar, _ := strconv.ParseInt(dataStr, 10, 32)
		val = reflect.ValueOf(int32(intVar))
	case reflect.Int64:
		val = reflect.ValueOf(data.(int64))
	case reflect.Uint:
		intVar, _ := strconv.Atoi(dataStr)
		val = reflect.ValueOf(uint(intVar))
	case reflect.Uint8:
		intVar, _ := strconv.ParseUint(dataStr, 10, 8)
		val = reflect.ValueOf(uint8(intVar))
	case reflect.Uint16:
		intVar, _ := strconv.ParseUint(dataStr, 10, 16)
		val = reflect.ValueOf(uint16(intVar))
	case reflect.Uint32:
		intVar, _ := strconv.ParseUint(dataStr, 10, 32)
		val = reflect.ValueOf(uint32(intVar))
	case reflect.Uint64:
		intVar, _ := strconv.ParseUint(dataStr, 10, 64)
		val = reflect.ValueOf(intVar)
	case reflect.Float32:
		floatVar, _ := strconv.ParseFloat(dataStr, 32)
		val = reflect.ValueOf(float32(floatVar))
	case reflect.Float64:
		floatVar, _ := strconv.ParseFloat(dataStr, 64)
		val = reflect.ValueOf(floatVar)
	default:
		err = errors.New("Unknown type %s", key)
	}

	// Check if the value is valid and not zero
	// IsZero panics if the value is invalid. Most functions and methods never return an invalid Value.
	isSet := val.IsValid() && !val.IsZero()

	if isSet {
		// Set the value for the field
		field.Set(val)
	}

	return
}

// toStr converts an interface{} type to a string representation.
// Supported types: bool, string, int (Int, Int8, Int16, Int32, Int64),
// uint (UInt, UInt8, UInt16, UInt32, UInt64), float (float32, float64)
//
// Parameters:
//   - data (interface{}): The value to be converted to a string.
//
// Returns:
//   - string: The string representation of the input value.
func toStr(data interface{}) (res string) {
	// Convert the input value to a string based on its type
	switch v := data.(type) {
	case float64:
		res = strconv.FormatFloat(data.(float64), 'f', 6, 64)
	case float32:
		res = strconv.FormatFloat(float64(data.(float32)), 'f', 6, 32)
	case int:
		res = strconv.FormatInt(int64(data.(int)), 10)
	case int64:
		res = strconv.FormatInt(data.(int64), 10)
	case uint:
		res = strconv.FormatUint(uint64(data.(uint)), 10)
	case uint64:
		res = strconv.FormatUint(data.(uint64), 10)
	case uint32:
		res = strconv.FormatUint(uint64(data.(uint32)), 10)
	case json.Number:
		res = data.(json.Number).String()
	case string:
		res = data.(string)
	case []byte:
		res = string(v)
	default:
		res = ""
	}

	return
}
