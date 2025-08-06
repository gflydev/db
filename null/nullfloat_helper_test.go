package null

import (
	"database/sql"
	"math"
	"testing"
)

func TestFloatNil(t *testing.T) {
	tests := []struct {
		name     string
		input    sql.NullFloat64
		expected *float64
	}{
		{
			name:     "valid positive float",
			input:    sql.NullFloat64{Float64: 3.14, Valid: true},
			expected: func() *float64 { f := 3.14; return &f }(),
		},
		{
			name:     "valid negative float",
			input:    sql.NullFloat64{Float64: -2.71, Valid: true},
			expected: func() *float64 { f := -2.71; return &f }(),
		},
		{
			name:     "valid zero",
			input:    sql.NullFloat64{Float64: 0.0, Valid: true},
			expected: func() *float64 { f := 0.0; return &f }(),
		},
		{
			name:     "valid infinity",
			input:    sql.NullFloat64{Float64: math.Inf(1), Valid: true},
			expected: func() *float64 { f := math.Inf(1); return &f }(),
		},
		{
			name:     "valid negative infinity",
			input:    sql.NullFloat64{Float64: math.Inf(-1), Valid: true},
			expected: func() *float64 { f := math.Inf(-1); return &f }(),
		},
		{
			name:     "valid NaN",
			input:    sql.NullFloat64{Float64: math.NaN(), Valid: true},
			expected: func() *float64 { f := math.NaN(); return &f }(),
		},
		{
			name:     "invalid null",
			input:    sql.NullFloat64{Float64: 0.0, Valid: false},
			expected: nil,
		},
		{
			name:     "invalid null with value",
			input:    sql.NullFloat64{Float64: 123.456, Valid: false},
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FloatNil(tt.input)
			if tt.expected == nil {
				if result != nil {
					t.Errorf("FloatNil(%v) = %v, want nil", tt.input, result)
				}
			} else {
				if result == nil {
					t.Errorf("FloatNil(%v) = nil, want %v", tt.input, *tt.expected)
				} else if math.IsNaN(*tt.expected) {
					if !math.IsNaN(*result) {
						t.Errorf("FloatNil(%v) = %v, want NaN", tt.input, *result)
					}
				} else if *result != *tt.expected {
					t.Errorf("FloatNil(%v) = %v, want %v", tt.input, *result, *tt.expected)
				}
			}
		})
	}
}

func TestFloatVal(t *testing.T) {
	tests := []struct {
		name     string
		input    sql.NullFloat64
		expected float64
	}{
		{
			name:     "valid float 3.14",
			input:    sql.NullFloat64{Float64: 3.14, Valid: true},
			expected: 3.14,
		},
		{
			name:     "valid float 0.0",
			input:    sql.NullFloat64{Float64: 0.0, Valid: true},
			expected: 0.0,
		},
		{
			name:     "valid negative float",
			input:    sql.NullFloat64{Float64: -123.456, Valid: true},
			expected: -123.456,
		},
		{
			name:     "valid max float64",
			input:    sql.NullFloat64{Float64: math.MaxFloat64, Valid: true},
			expected: math.MaxFloat64,
		},
		{
			name:     "valid min float64",
			input:    sql.NullFloat64{Float64: -math.MaxFloat64, Valid: true},
			expected: -math.MaxFloat64,
		},
		{
			name:     "invalid null",
			input:    sql.NullFloat64{Float64: 0.0, Valid: false},
			expected: 0.0,
		},
		{
			name:     "invalid null with value",
			input:    sql.NullFloat64{Float64: 123.456, Valid: false},
			expected: 0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FloatVal(tt.input)
			if result != tt.expected {
				t.Errorf("FloatVal(%v) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestFloat64(t *testing.T) {
	t.Run("float64 value positive", func(t *testing.T) {
		result := Float(3.14)
		expected := sql.NullFloat64{Float64: 3.14, Valid: true}
		if result != expected {
			t.Errorf("Float(3.14) = %v, want %v", result, expected)
		}
	})

	t.Run("float64 value negative", func(t *testing.T) {
		result := Float(-2.71)
		expected := sql.NullFloat64{Float64: -2.71, Valid: true}
		if result != expected {
			t.Errorf("Float(-2.71) = %v, want %v", result, expected)
		}
	})

	t.Run("float64 value zero", func(t *testing.T) {
		result := Float(0.0)
		expected := sql.NullFloat64{Float64: 0.0, Valid: true}
		if result != expected {
			t.Errorf("Float(0.0) = %v, want %v", result, expected)
		}
	})

	t.Run("float64 pointer positive", func(t *testing.T) {
		f := 1.23
		result := Float(&f)
		expected := sql.NullFloat64{Float64: 1.23, Valid: true}
		if result != expected {
			t.Errorf("Float(&1.23) = %v, want %v", result, expected)
		}
	})

	t.Run("float64 pointer negative", func(t *testing.T) {
		f := -4.56
		result := Float(&f)
		expected := sql.NullFloat64{Float64: -4.56, Valid: true}
		if result != expected {
			t.Errorf("Float(&-4.56) = %v, want %v", result, expected)
		}
	})

	t.Run("float64 pointer zero", func(t *testing.T) {
		f := 0.0
		result := Float(&f)
		expected := sql.NullFloat64{Float64: 0.0, Valid: true}
		if result != expected {
			t.Errorf("Float(&0.0) = %v, want %v", result, expected)
		}
	})

	t.Run("nil float64 pointer", func(t *testing.T) {
		result := Float((*float64)(nil))
		expected := sql.NullFloat64{Float64: 0, Valid: false}
		if result != expected {
			t.Errorf("Float((*float64)(nil)) = %v, want %v", result, expected)
		}
	})
}

func TestFloat64SpecialValues(t *testing.T) {
	t.Run("positive infinity", func(t *testing.T) {
		result := Float(math.Inf(1))
		expected := sql.NullFloat64{Float64: math.Inf(1), Valid: true}
		if result != expected {
			t.Errorf("Float(+Inf) = %v, want %v", result, expected)
		}
	})

	t.Run("negative infinity", func(t *testing.T) {
		result := Float(math.Inf(-1))
		expected := sql.NullFloat64{Float64: math.Inf(-1), Valid: true}
		if result != expected {
			t.Errorf("Float(-Inf) = %v, want %v", result, expected)
		}
	})

	t.Run("NaN", func(t *testing.T) {
		result := Float(math.NaN())
		if !math.IsNaN(result.Float64) || !result.Valid {
			t.Errorf("Float(NaN) = %v, want {NaN, true}", result)
		}
	})

	t.Run("max float64", func(t *testing.T) {
		result := Float(math.MaxFloat64)
		expected := sql.NullFloat64{Float64: math.MaxFloat64, Valid: true}
		if result != expected {
			t.Errorf("Float(MaxFloat64) = %v, want %v", result, expected)
		}
	})

	t.Run("smallest positive float64", func(t *testing.T) {
		result := Float(math.SmallestNonzeroFloat64)
		expected := sql.NullFloat64{Float64: math.SmallestNonzeroFloat64, Valid: true}
		if result != expected {
			t.Errorf("Float(SmallestNonzeroFloat64) = %v, want %v", result, expected)
		}
	})
}

func TestFloat64GenericTypeConstraints(t *testing.T) {
	// Test that the generic function works with both float64 and *float64 types
	t.Run("generic with float64", func(t *testing.T) {
		var f float64 = 42.42
		result := Float(f)
		expected := sql.NullFloat64{Float64: 42.42, Valid: true}
		if result != expected {
			t.Errorf("Float[float64](42.42) = %v, want %v", result, expected)
		}
	})

	t.Run("generic with *float64", func(t *testing.T) {
		f := 123.456
		var ptr *float64 = &f
		result := Float(ptr)
		expected := sql.NullFloat64{Float64: 123.456, Valid: true}
		if result != expected {
			t.Errorf("Float[*float64](&123.456) = %v, want %v", result, expected)
		}
	})
}

func TestFloat64EdgeCases(t *testing.T) {
	t.Run("very small positive number", func(t *testing.T) {
		small := 1e-100
		result := Float(small)
		expected := sql.NullFloat64{Float64: small, Valid: true}
		if result != expected {
			t.Errorf("Float(1e-100) = %v, want %v", result, expected)
		}
	})

	t.Run("very large positive number", func(t *testing.T) {
		large := 1e100
		result := Float(large)
		expected := sql.NullFloat64{Float64: large, Valid: true}
		if result != expected {
			t.Errorf("Float(1e100) = %v, want %v", result, expected)
		}
	})

	t.Run("negative zero", func(t *testing.T) {
		negZero := math.Copysign(0, -1)
		result := Float(negZero)
		expected := sql.NullFloat64{Float64: negZero, Valid: true}
		if result != expected {
			t.Errorf("Float(-0.0) = %v, want %v", result, expected)
		}
	})
}

// Benchmark tests

func BenchmarkFloatNil(b *testing.B) {
	nullFloat := sql.NullFloat64{Float64: 3.14, Valid: true}
	for i := 0; i < b.N; i++ {
		FloatNil(nullFloat)
	}
}

func BenchmarkFloat64(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Float(3.14)
	}
}

func BenchmarkFloat64Pointer(b *testing.B) {
	floatVal := 3.14
	for i := 0; i < b.N; i++ {
		Float(&floatVal)
	}
}
