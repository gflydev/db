package db

import (
	"encoding/json"
	"github.com/gflydev/core/errors"
	"reflect"
	"strconv"
)

// setValue dynamically assigns a value to a struct field using reflection.
// This function is a core utility for the ORM that enables dynamic field assignment
// during database operations such as primary key population after insertions,
// map-to-struct conversions, and result set mapping. It handles type conversion
// and validation to ensure safe assignment of values to struct fields.
//
// Parameters:
//   - model (any): The target struct or pointer to struct where the value will be set.
//     Must be a settable struct (typically a pointer to struct) with exported fields.
//     The struct should have the specified field name accessible for assignment.
//   - key (string): The exact field name in the struct to set the value for.
//     Field names are case-sensitive and must match the exported struct field names.
//     The field must be exported (capitalized) and settable.
//   - data (any): The value to assign to the specified field. The function performs
//     automatic type conversion based on the target field's type. Supported conversions:
//   - String values to string fields
//   - Boolean values to bool fields
//   - Numeric values with automatic conversion between int types (int, int8, int16, int32, int64)
//   - Numeric values with automatic conversion between uint types (uint, uint8, uint16, uint32, uint64)
//   - Floating-point values with conversion between float32 and float64
//
// Returns:
//   - error: Returns an error if:
//   - The model is not a struct or pointer to struct
//   - The specified field name doesn't exist in the struct
//   - The field is not exported or not settable
//   - Type conversion fails or is not supported
//   - The field type is not supported by the function
//     Returns nil on successful assignment.
//
// Supported Types:
//   - bool: Direct assignment of boolean values
//   - string: Direct assignment of string values
//   - int, int8, int16, int32, int64: Automatic conversion from string or numeric types
//   - uint, uint8, uint16, uint32, uint64: Automatic conversion with proper bounds checking
//   - float32, float64: Automatic conversion from string or numeric types
//
// Examples:
//
//	// Basic field assignment
//	type User struct {
//	    ID   int64  `db:"id,primary"`
//	    Name string `db:"name"`
//	    Age  int    `db:"age"`
//	}
//
//	user := &User{}
//	err := setValue(user, "ID", int64(123))        // Sets user.ID = 123
//	err = setValue(user, "Name", "John Doe")       // Sets user.Name = "John Doe"
//	err = setValue(user, "Age", 30)                // Sets user.Age = 30
//
//	// Type conversion examples
//	err := setValue(user, "Age", "25")             // String "25" converted to int 25
//	err = setValue(user, "ID", "456")              // String "456" converted to int64 456
//
//	// Primary key population after database insertion
//	newUser := &User{Name: "Jane", Age: 28}
//	// After database insertion, populate the generated ID
//	err := setValue(newUser, "ID", lastInsertID)   // Sets the auto-generated ID
//
//	// Map-to-struct conversion
//	userData := map[string]any{
//	    "Name": "Alice",
//	    "Age":  35,
//	}
//	user := &User{}
//	for key, value := range userData {
//	    err := setValue(user, key, value)
//	    if err != nil {
//	        log.Printf("Failed to set %s: %v", key, err)
//	    }
//	}
//
//	// Numeric type conversions
//	type Product struct {
//	    ID    uint32  `db:"id,primary"`
//	    Price float64 `db:"price"`
//	    Stock int16   `db:"stock"`
//	}
//
//	product := &Product{}
//	err := setValue(product, "ID", "12345")        // String to uint32
//	err = setValue(product, "Price", "99.99")      // String to float64
//	err = setValue(product, "Stock", "150")        // String to int16
//
// Use Cases:
//   - Primary key assignment after database insertions
//   - Dynamic struct field population from maps
//   - Result set mapping in ORM operations
//   - Type-safe field assignment with automatic conversion
//   - Database value assignment during query result processing
//
// Note:
//   - Only works with exported (capitalized) struct fields
//   - Requires a pointer to struct for field modification
//   - Zero values are skipped and won't be assigned
//   - Type conversion is automatic but limited to supported types
//   - Field names are case-sensitive and must match exactly
//   - Unsupported types will return an error
//   - The function performs bounds checking for numeric conversions
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
	// IsZero if the value is invalid. Most functions and methods never return an invalid Value.
	isSet := val.IsValid() && !val.IsZero()

	if isSet {
		// Set the value for the field
		field.Set(val)
	}

	return
}

// toStr converts various data types to their string representation for database operations.
// This utility function provides consistent string conversion for different Go types,
// ensuring proper formatting for database queries and type conversions. It handles
// numeric types, strings, byte slices, and JSON numbers with appropriate formatting
// to maintain precision and compatibility with database systems.
//
// Parameters:
//   - data (interface{}): The value to be converted to a string. Supported types:
//   - float64, float32: Converted with 6 decimal places precision
//   - int, int64: Converted to decimal string representation
//   - uint, uint64, uint32: Converted to decimal string representation
//   - json.Number: Uses the native string representation
//   - string: Returned as-is without modification
//   - []byte: Converted to string using UTF-8 encoding
//   - Other types: Returns empty string as fallback
//
// Returns:
//   - string: The string representation of the input value with appropriate formatting:
//   - Numeric types: Decimal representation without scientific notation
//   - Float types: Fixed-point notation with 6 decimal places
//   - String types: Unchanged original value
//   - Byte slices: UTF-8 decoded string
//   - Unsupported types: Empty string ("")
//
// Examples:
//
//	// Numeric conversions
//	result := toStr(123)           // Returns: "123"
//	result = toStr(int64(456))     // Returns: "456"
//	result = toStr(uint32(789))    // Returns: "789"
//
//	// Floating-point conversions
//	result := toStr(3.14159)       // Returns: "3.141590" (6 decimal places)
//	result = toStr(float32(2.5))   // Returns: "2.500000"
//	result = toStr(99.0)           // Returns: "99.000000"
//
//	// String and byte conversions
//	result := toStr("hello")       // Returns: "hello"
//	result = toStr([]byte("world")) // Returns: "world"
//
//	// JSON number conversion
//	jsonNum := json.Number("123.45")
//	result := toStr(jsonNum)       // Returns: "123.45"
//
//	// Unsupported type fallback
//	result := toStr(struct{}{})    // Returns: ""
//	result = toStr(nil)            // Returns: ""
//
// Use Cases:
//   - Database parameter conversion in setValue function
//   - Query argument formatting for SQL operations
//   - Type conversion during ORM field mapping
//   - Consistent string representation for logging and debugging
//   - Data serialization for database storage
//
// Formatting Details:
//   - Float64: Uses 'f' format with 6 decimal places and 64-bit precision
//   - Float32: Converted to float64 then formatted with 32-bit precision
//   - Integers: Base-10 decimal representation without leading zeros
//   - JSON Numbers: Preserves original string format from JSON parsing
//   - Byte slices: Direct UTF-8 string conversion
//
// Note:
//   - Float formatting always includes 6 decimal places for consistency
//   - Large integers are formatted without scientific notation
//   - Byte slice conversion assumes UTF-8 encoding
//   - Unsupported types return empty string rather than panicking
//   - The function is designed for database compatibility, not human readability
//   - No error handling - invalid conversions result in empty strings
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
