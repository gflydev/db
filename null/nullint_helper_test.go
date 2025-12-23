package null

import (
	"database/sql"
	"math"
	"testing"
)

// Tests for Int64 functions

func TestInt64Nil(t *testing.T) {
	tests := []struct {
		name     string
		input    sql.NullInt64
		expected *int64
	}{
		{
			name:     "valid positive int64",
			input:    sql.NullInt64{Int64: 123456789, Valid: true},
			expected: func() *int64 { i := int64(123456789); return &i }(),
		},
		{
			name:     "valid negative int64",
			input:    sql.NullInt64{Int64: -987654321, Valid: true},
			expected: func() *int64 { i := int64(-987654321); return &i }(),
		},
		{
			name:     "valid zero",
			input:    sql.NullInt64{Int64: 0, Valid: true},
			expected: func() *int64 { i := int64(0); return &i }(),
		},
		{
			name:     "valid max int64",
			input:    sql.NullInt64{Int64: math.MaxInt64, Valid: true},
			expected: func() *int64 { i := int64(math.MaxInt64); return &i }(),
		},
		{
			name:     "valid min int64",
			input:    sql.NullInt64{Int64: math.MinInt64, Valid: true},
			expected: func() *int64 { i := int64(math.MinInt64); return &i }(),
		},
		{
			name:     "invalid null",
			input:    sql.NullInt64{Int64: 0, Valid: false},
			expected: nil,
		},
		{
			name:     "invalid null with value",
			input:    sql.NullInt64{Int64: 12345, Valid: false},
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Int64Nil(tt.input)
			if tt.expected == nil {
				if result != nil {
					t.Errorf("Int64Nil(%v) = %v, want nil", tt.input, result)
				}
			} else {
				if result == nil {
					t.Errorf("Int64Nil(%v) = nil, want %v", tt.input, *tt.expected)
				} else if *result != *tt.expected {
					t.Errorf("Int64Nil(%v) = %v, want %v", tt.input, *result, *tt.expected)
				}
			}
		})
	}
}

func TestInt64Val(t *testing.T) {
	tests := []struct {
		name     string
		input    sql.NullInt64
		expected int64
	}{
		{
			name:     "valid int64 42",
			input:    sql.NullInt64{Int64: 42, Valid: true},
			expected: int64(42),
		},
		{
			name:     "valid int64 0",
			input:    sql.NullInt64{Int64: 0, Valid: true},
			expected: int64(0),
		},
		{
			name:     "valid negative int64",
			input:    sql.NullInt64{Int64: -123, Valid: true},
			expected: int64(-123),
		},
		{
			name:     "valid max int64",
			input:    sql.NullInt64{Int64: math.MaxInt64, Valid: true},
			expected: int64(math.MaxInt64),
		},
		{
			name:     "valid min int64",
			input:    sql.NullInt64{Int64: math.MinInt64, Valid: true},
			expected: int64(math.MinInt64),
		},
		{
			name:     "invalid null",
			input:    sql.NullInt64{Int64: 0, Valid: false},
			expected: int64(0),
		},
		{
			name:     "invalid null with value",
			input:    sql.NullInt64{Int64: 123, Valid: false},
			expected: int64(0),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Int64Val(tt.input)
			if result != tt.expected {
				t.Errorf("Int64Val(%v) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestInt64(t *testing.T) {
	t.Run("int64 value positive", func(t *testing.T) {
		result := Int64(int64(123456789))
		expected := sql.NullInt64{Int64: 123456789, Valid: true}
		if result != expected {
			t.Errorf("Int64(123456789) = %v, want %v", result, expected)
		}
	})

	t.Run("int64 value negative", func(t *testing.T) {
		result := Int64(int64(-987654321))
		expected := sql.NullInt64{Int64: -987654321, Valid: true}
		if result != expected {
			t.Errorf("Int64(-987654321) = %v, want %v", result, expected)
		}
	})

	t.Run("int64 value zero", func(t *testing.T) {
		result := Int64(int64(0))
		expected := sql.NullInt64{Int64: 0, Valid: true}
		if result != expected {
			t.Errorf("Int64(0) = %v, want %v", result, expected)
		}
	})

	t.Run("int64 pointer positive", func(t *testing.T) {
		i := int64(42)
		result := Int64(&i)
		expected := sql.NullInt64{Int64: 42, Valid: true}
		if result != expected {
			t.Errorf("Int64(&42) = %v, want %v", result, expected)
		}
	})

	t.Run("int64 pointer negative", func(t *testing.T) {
		i := int64(-42)
		result := Int64(&i)
		expected := sql.NullInt64{Int64: -42, Valid: true}
		if result != expected {
			t.Errorf("Int64(&-42) = %v, want %v", result, expected)
		}
	})

	t.Run("nil int64 pointer", func(t *testing.T) {
		result := Int64((*int64)(nil))
		expected := sql.NullInt64{Int64: 0, Valid: false}
		if result != expected {
			t.Errorf("Int64((*int64)(nil)) = %v, want %v", result, expected)
		}
	})

	t.Run("max int64", func(t *testing.T) {
		result := Int64(int64(math.MaxInt64))
		expected := sql.NullInt64{Int64: math.MaxInt64, Valid: true}
		if result != expected {
			t.Errorf("Int64(MaxInt64) = %v, want %v", result, expected)
		}
	})

	t.Run("min int64", func(t *testing.T) {
		result := Int64(int64(math.MinInt64))
		expected := sql.NullInt64{Int64: math.MinInt64, Valid: true}
		if result != expected {
			t.Errorf("Int64(MinInt64) = %v, want %v", result, expected)
		}
	})
}

// Tests for Int32 functions

func TestInt32Nil(t *testing.T) {
	tests := []struct {
		name     string
		input    sql.NullInt32
		expected *int32
	}{
		{
			name:     "valid positive int32",
			input:    sql.NullInt32{Int32: 123456, Valid: true},
			expected: func() *int32 { i := int32(123456); return &i }(),
		},
		{
			name:     "valid negative int32",
			input:    sql.NullInt32{Int32: -654321, Valid: true},
			expected: func() *int32 { i := int32(-654321); return &i }(),
		},
		{
			name:     "valid zero",
			input:    sql.NullInt32{Int32: 0, Valid: true},
			expected: func() *int32 { i := int32(0); return &i }(),
		},
		{
			name:     "invalid null",
			input:    sql.NullInt32{Int32: 0, Valid: false},
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Int32Nil(tt.input)
			if tt.expected == nil {
				if result != nil {
					t.Errorf("Int32Nil(%v) = %v, want nil", tt.input, result)
				}
			} else {
				if result == nil {
					t.Errorf("Int32Nil(%v) = nil, want %v", tt.input, *tt.expected)
				} else if *result != *tt.expected {
					t.Errorf("Int32Nil(%v) = %v, want %v", tt.input, *result, *tt.expected)
				}
			}
		})
	}
}

func TestInt32Val(t *testing.T) {
	tests := []struct {
		name     string
		input    sql.NullInt32
		expected int32
	}{
		{
			name:     "valid int32 42",
			input:    sql.NullInt32{Int32: 42, Valid: true},
			expected: int32(42),
		},
		{
			name:     "valid int32 0",
			input:    sql.NullInt32{Int32: 0, Valid: true},
			expected: int32(0),
		},
		{
			name:     "valid negative int32",
			input:    sql.NullInt32{Int32: -123, Valid: true},
			expected: int32(-123),
		},
		{
			name:     "valid max int32",
			input:    sql.NullInt32{Int32: math.MaxInt32, Valid: true},
			expected: int32(math.MaxInt32),
		},
		{
			name:     "valid min int32",
			input:    sql.NullInt32{Int32: math.MinInt32, Valid: true},
			expected: int32(math.MinInt32),
		},
		{
			name:     "invalid null",
			input:    sql.NullInt32{Int32: 0, Valid: false},
			expected: int32(0),
		},
		{
			name:     "invalid null with value",
			input:    sql.NullInt32{Int32: 123, Valid: false},
			expected: int32(0),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Int32Val(tt.input)
			if result != tt.expected {
				t.Errorf("Int32Val(%v) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestInt32(t *testing.T) {
	t.Run("int32 value positive", func(t *testing.T) {
		result := Int32(int32(123456))
		expected := sql.NullInt32{Int32: 123456, Valid: true}
		if result != expected {
			t.Errorf("Int32(123456) = %v, want %v", result, expected)
		}
	})

	t.Run("int32 value negative", func(t *testing.T) {
		result := Int32(int32(-654321))
		expected := sql.NullInt32{Int32: -654321, Valid: true}
		if result != expected {
			t.Errorf("Int32(-654321) = %v, want %v", result, expected)
		}
	})

	t.Run("int32 pointer", func(t *testing.T) {
		i := int32(42)
		result := Int32(&i)
		expected := sql.NullInt32{Int32: 42, Valid: true}
		if result != expected {
			t.Errorf("Int32(&42) = %v, want %v", result, expected)
		}
	})

	t.Run("nil int32 pointer", func(t *testing.T) {
		result := Int32((*int32)(nil))
		expected := sql.NullInt32{Int32: 0, Valid: false}
		if result != expected {
			t.Errorf("Int32((*int32)(nil)) = %v, want %v", result, expected)
		}
	})

	t.Run("max int32", func(t *testing.T) {
		result := Int32(int32(math.MaxInt32))
		expected := sql.NullInt32{Int32: math.MaxInt32, Valid: true}
		if result != expected {
			t.Errorf("Int32(MaxInt32) = %v, want %v", result, expected)
		}
	})

	t.Run("min int32", func(t *testing.T) {
		result := Int32(int32(math.MinInt32))
		expected := sql.NullInt32{Int32: math.MinInt32, Valid: true}
		if result != expected {
			t.Errorf("Int32(MinInt32) = %v, want %v", result, expected)
		}
	})
}

// Tests for Int16 functions

func TestInt16Nil(t *testing.T) {
	tests := []struct {
		name     string
		input    sql.NullInt16
		expected *int16
	}{
		{
			name:     "valid positive int16",
			input:    sql.NullInt16{Int16: 12345, Valid: true},
			expected: func() *int16 { i := int16(12345); return &i }(),
		},
		{
			name:     "valid negative int16",
			input:    sql.NullInt16{Int16: -6543, Valid: true},
			expected: func() *int16 { i := int16(-6543); return &i }(),
		},
		{
			name:     "valid zero",
			input:    sql.NullInt16{Int16: 0, Valid: true},
			expected: func() *int16 { i := int16(0); return &i }(),
		},
		{
			name:     "invalid null",
			input:    sql.NullInt16{Int16: 0, Valid: false},
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Int16Nil(tt.input)
			if tt.expected == nil {
				if result != nil {
					t.Errorf("Int16Nil(%v) = %v, want nil", tt.input, result)
				}
			} else {
				if result == nil {
					t.Errorf("Int16Nil(%v) = nil, want %v", tt.input, *tt.expected)
				} else if *result != *tt.expected {
					t.Errorf("Int16Nil(%v) = %v, want %v", tt.input, *result, *tt.expected)
				}
			}
		})
	}
}

func TestInt16Val(t *testing.T) {
	tests := []struct {
		name     string
		input    sql.NullInt16
		expected int16
	}{
		{
			name:     "valid int16 42",
			input:    sql.NullInt16{Int16: 42, Valid: true},
			expected: int16(42),
		},
		{
			name:     "valid int16 0",
			input:    sql.NullInt16{Int16: 0, Valid: true},
			expected: int16(0),
		},
		{
			name:     "valid negative int16",
			input:    sql.NullInt16{Int16: -123, Valid: true},
			expected: int16(-123),
		},
		{
			name:     "valid max int16",
			input:    sql.NullInt16{Int16: math.MaxInt16, Valid: true},
			expected: int16(math.MaxInt16),
		},
		{
			name:     "valid min int16",
			input:    sql.NullInt16{Int16: math.MinInt16, Valid: true},
			expected: int16(math.MinInt16),
		},
		{
			name:     "invalid null",
			input:    sql.NullInt16{Int16: 0, Valid: false},
			expected: int16(0),
		},
		{
			name:     "invalid null with value",
			input:    sql.NullInt16{Int16: 123, Valid: false},
			expected: int16(0),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Int16Val(tt.input)
			if result != tt.expected {
				t.Errorf("Int16Val(%v) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestInt16(t *testing.T) {
	t.Run("int16 value positive", func(t *testing.T) {
		result := Int16(int16(12345))
		expected := sql.NullInt16{Int16: 12345, Valid: true}
		if result != expected {
			t.Errorf("Int16(12345) = %v, want %v", result, expected)
		}
	})

	t.Run("int16 value negative", func(t *testing.T) {
		result := Int16(int16(-6543))
		expected := sql.NullInt16{Int16: -6543, Valid: true}
		if result != expected {
			t.Errorf("Int16(-6543) = %v, want %v", result, expected)
		}
	})

	t.Run("int16 pointer", func(t *testing.T) {
		i := int16(42)
		result := Int16(&i)
		expected := sql.NullInt16{Int16: 42, Valid: true}
		if result != expected {
			t.Errorf("Int16(&42) = %v, want %v", result, expected)
		}
	})

	t.Run("nil int16 pointer", func(t *testing.T) {
		result := Int16((*int16)(nil))
		expected := sql.NullInt16{Int16: 0, Valid: false}
		if result != expected {
			t.Errorf("Int16((*int16)(nil)) = %v, want %v", result, expected)
		}
	})

	t.Run("max int16", func(t *testing.T) {
		result := Int16(int16(math.MaxInt16))
		expected := sql.NullInt16{Int16: math.MaxInt16, Valid: true}
		if result != expected {
			t.Errorf("Int16(MaxInt16) = %v, want %v", result, expected)
		}
	})

	t.Run("min int16", func(t *testing.T) {
		result := Int16(int16(math.MinInt16))
		expected := sql.NullInt16{Int16: math.MinInt16, Valid: true}
		if result != expected {
			t.Errorf("Int16(MinInt16) = %v, want %v", result, expected)
		}
	})
}

// Generic type constraint tests
func TestIntGenericTypeConstraints(t *testing.T) {
	t.Run("int64 generic constraints", func(t *testing.T) {
		var i int64 = 42
		result := Int64(i)
		expected := sql.NullInt64{Int64: 42, Valid: true}
		if result != expected {
			t.Errorf("Int64[int64](42) = %v, want %v", result, expected)
		}

		var ptr *int64 = &i
		result = Int64(ptr)
		if result != expected {
			t.Errorf("Int64[*int64](&42) = %v, want %v", result, expected)
		}
	})

	t.Run("int32 generic constraints", func(t *testing.T) {
		var i int32 = 42
		result := Int32(i)
		expected := sql.NullInt32{Int32: 42, Valid: true}
		if result != expected {
			t.Errorf("Int32[int32](42) = %v, want %v", result, expected)
		}

		var ptr *int32 = &i
		result = Int32(ptr)
		if result != expected {
			t.Errorf("Int32[*int32](&42) = %v, want %v", result, expected)
		}
	})

	t.Run("int16 generic constraints", func(t *testing.T) {
		var i int16 = 42
		result := Int16(i)
		expected := sql.NullInt16{Int16: 42, Valid: true}
		if result != expected {
			t.Errorf("Int16[int16](42) = %v, want %v", result, expected)
		}

		var ptr *int16 = &i
		result = Int16(ptr)
		if result != expected {
			t.Errorf("Int16[*int16](&42) = %v, want %v", result, expected)
		}
	})
}

// Benchmark tests

func BenchmarkInt64Nil(b *testing.B) {
	nullInt := sql.NullInt64{Int64: 123456789, Valid: true}
	for i := 0; i < b.N; i++ {
		Int64Nil(nullInt)
	}
}

func BenchmarkInt64(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Int64(int64(123456789))
	}
}

func BenchmarkInt32Nil(b *testing.B) {
	nullInt := sql.NullInt32{Int32: 123456, Valid: true}
	for i := 0; i < b.N; i++ {
		Int32Nil(nullInt)
	}
}

func BenchmarkInt32(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Int32(int32(123456))
	}
}

func BenchmarkInt16Nil(b *testing.B) {
	nullInt := sql.NullInt16{Int16: 12345, Valid: true}
	for i := 0; i < b.N; i++ {
		Int16Nil(nullInt)
	}
}

func BenchmarkInt16(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Int16(int16(12345))
	}
}

// Tests for Int64NilInt, Int32NilInt, Int16NilInt functions

func TestInt64NilInt(t *testing.T) {
	tests := []struct {
		name     string
		input    sql.NullInt64
		expected *int
	}{
		{
			name:     "valid positive int64",
			input:    sql.NullInt64{Int64: 123456789, Valid: true},
			expected: func() *int { i := int(123456789); return &i }(),
		},
		{
			name:     "valid negative int64",
			input:    sql.NullInt64{Int64: -987654321, Valid: true},
			expected: func() *int { i := int(-987654321); return &i }(),
		},
		{
			name:     "valid zero",
			input:    sql.NullInt64{Int64: 0, Valid: true},
			expected: func() *int { i := int(0); return &i }(),
		},
		{
			name:     "valid max int64 (may truncate on 32-bit systems)",
			input:    sql.NullInt64{Int64: math.MaxInt64, Valid: true},
			expected: func() *int { i := int(math.MaxInt64); return &i }(),
		},
		{
			name:     "valid min int64 (may truncate on 32-bit systems)",
			input:    sql.NullInt64{Int64: math.MinInt64, Valid: true},
			expected: func() *int { i := int(math.MinInt64); return &i }(),
		},
		{
			name:     "invalid null",
			input:    sql.NullInt64{Int64: 0, Valid: false},
			expected: nil,
		},
		{
			name:     "invalid null with value",
			input:    sql.NullInt64{Int64: 12345, Valid: false},
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Int64NilInt(tt.input)
			if tt.expected == nil {
				if result != nil {
					t.Errorf("Int64NilInt(%v) = %v, want nil", tt.input, result)
				}
			} else {
				if result == nil {
					t.Errorf("Int64NilInt(%v) = nil, want %v", tt.input, *tt.expected)
				} else if *result != *tt.expected {
					t.Errorf("Int64NilInt(%v) = %v, want %v", tt.input, *result, *tt.expected)
				}
			}
		})
	}
}

func TestInt32NilInt(t *testing.T) {
	tests := []struct {
		name     string
		input    sql.NullInt32
		expected *int
	}{
		{
			name:     "valid positive int32",
			input:    sql.NullInt32{Int32: 123456, Valid: true},
			expected: func() *int { i := int(123456); return &i }(),
		},
		{
			name:     "valid negative int32",
			input:    sql.NullInt32{Int32: -654321, Valid: true},
			expected: func() *int { i := int(-654321); return &i }(),
		},
		{
			name:     "valid zero",
			input:    sql.NullInt32{Int32: 0, Valid: true},
			expected: func() *int { i := int(0); return &i }(),
		},
		{
			name:     "valid max int32",
			input:    sql.NullInt32{Int32: math.MaxInt32, Valid: true},
			expected: func() *int { i := int(math.MaxInt32); return &i }(),
		},
		{
			name:     "valid min int32",
			input:    sql.NullInt32{Int32: math.MinInt32, Valid: true},
			expected: func() *int { i := int(math.MinInt32); return &i }(),
		},
		{
			name:     "invalid null",
			input:    sql.NullInt32{Int32: 0, Valid: false},
			expected: nil,
		},
		{
			name:     "invalid null with value",
			input:    sql.NullInt32{Int32: 12345, Valid: false},
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Int32NilInt(tt.input)
			if tt.expected == nil {
				if result != nil {
					t.Errorf("Int32NilInt(%v) = %v, want nil", tt.input, result)
				}
			} else {
				if result == nil {
					t.Errorf("Int32NilInt(%v) = nil, want %v", tt.input, *tt.expected)
				} else if *result != *tt.expected {
					t.Errorf("Int32NilInt(%v) = %v, want %v", tt.input, *result, *tt.expected)
				}
			}
		})
	}
}

func TestInt16NilInt(t *testing.T) {
	tests := []struct {
		name     string
		input    sql.NullInt16
		expected *int
	}{
		{
			name:     "valid positive int16",
			input:    sql.NullInt16{Int16: 12345, Valid: true},
			expected: func() *int { i := int(12345); return &i }(),
		},
		{
			name:     "valid negative int16",
			input:    sql.NullInt16{Int16: -6543, Valid: true},
			expected: func() *int { i := int(-6543); return &i }(),
		},
		{
			name:     "valid zero",
			input:    sql.NullInt16{Int16: 0, Valid: true},
			expected: func() *int { i := int(0); return &i }(),
		},
		{
			name:     "valid max int16",
			input:    sql.NullInt16{Int16: math.MaxInt16, Valid: true},
			expected: func() *int { i := int(math.MaxInt16); return &i }(),
		},
		{
			name:     "valid min int16",
			input:    sql.NullInt16{Int16: math.MinInt16, Valid: true},
			expected: func() *int { i := int(math.MinInt16); return &i }(),
		},
		{
			name:     "invalid null",
			input:    sql.NullInt16{Int16: 0, Valid: false},
			expected: nil,
		},
		{
			name:     "invalid null with value",
			input:    sql.NullInt16{Int16: 12345, Valid: false},
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Int16NilInt(tt.input)
			if tt.expected == nil {
				if result != nil {
					t.Errorf("Int16NilInt(%v) = %v, want nil", tt.input, result)
				}
			} else {
				if result == nil {
					t.Errorf("Int16NilInt(%v) = nil, want %v", tt.input, *tt.expected)
				} else if *result != *tt.expected {
					t.Errorf("Int16NilInt(%v) = %v, want %v", tt.input, *result, *tt.expected)
				}
			}
		})
	}
}

func BenchmarkInt64NilInt(b *testing.B) {
	nullInt := sql.NullInt64{Int64: 123456789, Valid: true}
	for i := 0; i < b.N; i++ {
		Int64NilInt(nullInt)
	}
}

func BenchmarkInt32NilInt(b *testing.B) {
	nullInt := sql.NullInt32{Int32: 123456, Valid: true}
	for i := 0; i < b.N; i++ {
		Int32NilInt(nullInt)
	}
}

func BenchmarkInt16NilInt(b *testing.B) {
	nullInt := sql.NullInt16{Int16: 12345, Valid: true}
	for i := 0; i < b.N; i++ {
		Int16NilInt(nullInt)
	}
}
