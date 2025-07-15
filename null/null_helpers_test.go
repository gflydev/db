package null

import (
	"testing"
	"time"
)

func TestTime(t *testing.T) {
	// Test with valid time.Time
	now := time.Now()
	result := Time(now)
	if !result.Valid {
		t.Error("Expected Valid to be true for time.Time input")
	}
	if result.Time != now {
		t.Error("Expected Time value to match input")
	}

	// Test with valid *time.Time
	timePtr := &now
	result = Time(timePtr)
	if !result.Valid {
		t.Error("Expected Valid to be true for *time.Time input")
	}
	if result.Time != now {
		t.Error("Expected Time value to match dereferenced input")
	}

	// Test with nil *time.Time
	var nilTimePtr *time.Time
	result = Time(nilTimePtr)
	if result.Valid {
		t.Error("Expected Valid to be false for nil *time.Time input")
	}
	if !result.Time.IsZero() {
		t.Error("Expected Time to be zero value for nil input")
	}

	// Test with invalid type
	result = Time("invalid")
	if result.Valid {
		t.Error("Expected Valid to be false for invalid type")
	}
	if !result.Time.IsZero() {
		t.Error("Expected Time to be zero value for invalid type")
	}
}

func TestString(t *testing.T) {
	// Test with valid string
	testStr := "test string"
	result := String(testStr)
	if !result.Valid {
		t.Error("Expected Valid to be true for string input")
	}
	if result.String != testStr {
		t.Error("Expected String value to match input")
	}

	// Test with valid *string
	strPtr := &testStr
	result = String(strPtr)
	if !result.Valid {
		t.Error("Expected Valid to be true for *string input")
	}
	if result.String != testStr {
		t.Error("Expected String value to match dereferenced input")
	}

	// Test with nil *string
	var nilStrPtr *string
	result = String(nilStrPtr)
	if result.Valid {
		t.Error("Expected Valid to be false for nil *string input")
	}
	if result.String != "" {
		t.Error("Expected String to be empty for nil input")
	}

	// Test with invalid type
	result = String(123)
	if result.Valid {
		t.Error("Expected Valid to be false for invalid type")
	}
	if result.String != "" {
		t.Error("Expected String to be empty for invalid type")
	}
}

func TestInt64(t *testing.T) {
	// Test with valid int64
	var testInt int64 = 12345
	result := Int64(testInt)
	if !result.Valid {
		t.Error("Expected Valid to be true for int64 input")
	}
	if result.Int64 != testInt {
		t.Error("Expected Int64 value to match input")
	}

	// Test with valid *int64
	intPtr := &testInt
	result = Int64(intPtr)
	if !result.Valid {
		t.Error("Expected Valid to be true for *int64 input")
	}
	if result.Int64 != testInt {
		t.Error("Expected Int64 value to match dereferenced input")
	}

	// Test with nil *int64
	var nilIntPtr *int64
	result = Int64(nilIntPtr)
	if result.Valid {
		t.Error("Expected Valid to be false for nil *int64 input")
	}
	if result.Int64 != 0 {
		t.Error("Expected Int64 to be 0 for nil input")
	}

	// Test with invalid type
	result = Int64("invalid")
	if result.Valid {
		t.Error("Expected Valid to be false for invalid type")
	}
	if result.Int64 != 0 {
		t.Error("Expected Int64 to be 0 for invalid type")
	}
}

func TestInt32(t *testing.T) {
	// Test with valid int32
	var testInt int32 = 12345
	result := Int32(testInt)
	if !result.Valid {
		t.Error("Expected Valid to be true for int32 input")
	}
	if result.Int32 != testInt {
		t.Error("Expected Int32 value to match input")
	}

	// Test with valid *int32
	intPtr := &testInt
	result = Int32(intPtr)
	if !result.Valid {
		t.Error("Expected Valid to be true for *int32 input")
	}
	if result.Int32 != testInt {
		t.Error("Expected Int32 value to match dereferenced input")
	}

	// Test with nil *int32
	var nilIntPtr *int32
	result = Int32(nilIntPtr)
	if result.Valid {
		t.Error("Expected Valid to be false for nil *int32 input")
	}
	if result.Int32 != 0 {
		t.Error("Expected Int32 to be 0 for nil input")
	}

	// Test with invalid type
	result = Int32("invalid")
	if result.Valid {
		t.Error("Expected Valid to be false for invalid type")
	}
	if result.Int32 != 0 {
		t.Error("Expected Int32 to be 0 for invalid type")
	}
}

func TestInt16(t *testing.T) {
	// Test with valid int16
	var testInt int16 = 12345
	result := Int16(testInt)
	if !result.Valid {
		t.Error("Expected Valid to be true for int16 input")
	}
	if result.Int16 != testInt {
		t.Error("Expected Int16 value to match input")
	}

	// Test with valid *int16
	intPtr := &testInt
	result = Int16(intPtr)
	if !result.Valid {
		t.Error("Expected Valid to be true for *int16 input")
	}
	if result.Int16 != testInt {
		t.Error("Expected Int16 value to match dereferenced input")
	}

	// Test with nil *int16
	var nilIntPtr *int16
	result = Int16(nilIntPtr)
	if result.Valid {
		t.Error("Expected Valid to be false for nil *int16 input")
	}
	if result.Int16 != 0 {
		t.Error("Expected Int16 to be 0 for nil input")
	}

	// Test with invalid type
	result = Int16("invalid")
	if result.Valid {
		t.Error("Expected Valid to be false for invalid type")
	}
	if result.Int16 != 0 {
		t.Error("Expected Int16 to be 0 for invalid type")
	}
}

func TestFloat64(t *testing.T) {
	// Test with valid float64
	var testFloat float64 = 123.45
	result := Float64(testFloat)
	if !result.Valid {
		t.Error("Expected Valid to be true for float64 input")
	}
	if result.Float64 != testFloat {
		t.Error("Expected Float64 value to match input")
	}

	// Test with valid *float64
	floatPtr := &testFloat
	result = Float64(floatPtr)
	if !result.Valid {
		t.Error("Expected Valid to be true for *float64 input")
	}
	if result.Float64 != testFloat {
		t.Error("Expected Float64 value to match dereferenced input")
	}

	// Test with nil *float64
	var nilFloatPtr *float64
	result = Float64(nilFloatPtr)
	if result.Valid {
		t.Error("Expected Valid to be false for nil *float64 input")
	}
	if result.Float64 != 0 {
		t.Error("Expected Float64 to be 0 for nil input")
	}

	// Test with invalid type
	result = Float64("invalid")
	if result.Valid {
		t.Error("Expected Valid to be false for invalid type")
	}
	if result.Float64 != 0 {
		t.Error("Expected Float64 to be 0 for invalid type")
	}
}

func TestByte(t *testing.T) {
	// Test with valid byte
	var testByte byte = 255
	result := Byte(testByte)
	if !result.Valid {
		t.Error("Expected Valid to be true for byte input")
	}
	if result.Byte != testByte {
		t.Error("Expected Byte value to match input")
	}

	// Test with valid *byte
	bytePtr := &testByte
	result = Byte(bytePtr)
	if !result.Valid {
		t.Error("Expected Valid to be true for *byte input")
	}
	if result.Byte != testByte {
		t.Error("Expected Byte value to match dereferenced input")
	}

	// Test with nil *byte
	var nilBytePtr *byte
	result = Byte(nilBytePtr)
	if result.Valid {
		t.Error("Expected Valid to be false for nil *byte input")
	}
	if result.Byte != 0 {
		t.Error("Expected Byte to be 0 for nil input")
	}

	// Test with invalid type
	result = Byte("invalid")
	if result.Valid {
		t.Error("Expected Valid to be false for invalid type")
	}
	if result.Byte != 0 {
		t.Error("Expected Byte to be 0 for invalid type")
	}
}

func TestBool(t *testing.T) {
	// Test with valid bool (true)
	testBool := true
	result := Bool(testBool)
	if !result.Valid {
		t.Error("Expected Valid to be true for bool input")
	}
	if result.Bool != testBool {
		t.Error("Expected Bool value to match input")
	}

	// Test with valid bool (false)
	testBool = false
	result = Bool(testBool)
	if !result.Valid {
		t.Error("Expected Valid to be true for bool input")
	}
	if result.Bool != testBool {
		t.Error("Expected Bool value to match input")
	}

	// Test with valid *bool
	boolPtr := &testBool
	result = Bool(boolPtr)
	if !result.Valid {
		t.Error("Expected Valid to be true for *bool input")
	}
	if result.Bool != testBool {
		t.Error("Expected Bool value to match dereferenced input")
	}

	// Test with nil *bool
	var nilBoolPtr *bool
	result = Bool(nilBoolPtr)
	if result.Valid {
		t.Error("Expected Valid to be false for nil *bool input")
	}
	if result.Bool != false {
		t.Error("Expected Bool to be false for nil input")
	}

	// Test with invalid type
	result = Bool("invalid")
	if result.Valid {
		t.Error("Expected Valid to be false for invalid type")
	}
	if result.Bool != false {
		t.Error("Expected Bool to be false for invalid type")
	}
}
