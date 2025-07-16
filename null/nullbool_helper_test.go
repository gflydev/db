package null

import (
	"database/sql"
	"database/sql/driver"
	"testing"
)

func TestBoolAny(t *testing.T) {
	tests := []struct {
		name     string
		input    sql.NullBool
		expected driver.Value
	}{
		{
			name:     "valid true",
			input:    sql.NullBool{Bool: true, Valid: true},
			expected: true,
		},
		{
			name:     "valid false",
			input:    sql.NullBool{Bool: false, Valid: true},
			expected: false,
		},
		{
			name:     "invalid null",
			input:    sql.NullBool{Bool: false, Valid: false},
			expected: nil,
		},
		{
			name:     "invalid null with true value",
			input:    sql.NullBool{Bool: true, Valid: false},
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := BoolAny(tt.input)
			if result != tt.expected {
				t.Errorf("BoolAny(%v) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestBoolNil(t *testing.T) {
	tests := []struct {
		name     string
		input    sql.NullBool
		expected *bool
	}{
		{
			name:     "valid true",
			input:    sql.NullBool{Bool: true, Valid: true},
			expected: func() *bool { b := true; return &b }(),
		},
		{
			name:     "valid false",
			input:    sql.NullBool{Bool: false, Valid: true},
			expected: func() *bool { b := false; return &b }(),
		},
		{
			name:     "invalid null",
			input:    sql.NullBool{Bool: false, Valid: false},
			expected: nil,
		},
		{
			name:     "invalid null with true value",
			input:    sql.NullBool{Bool: true, Valid: false},
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := BoolNil(tt.input)
			if tt.expected == nil {
				if result != nil {
					t.Errorf("BoolNil(%v) = %v, want nil", tt.input, result)
				}
			} else {
				if result == nil {
					t.Errorf("BoolNil(%v) = nil, want %v", tt.input, *tt.expected)
				} else if *result != *tt.expected {
					t.Errorf("BoolNil(%v) = %v, want %v", tt.input, *result, *tt.expected)
				}
			}
		})
	}
}

func TestBool(t *testing.T) {
	t.Run("bool value true", func(t *testing.T) {
		result := Bool(true)
		expected := sql.NullBool{Bool: true, Valid: true}
		if result != expected {
			t.Errorf("Bool(true) = %v, want %v", result, expected)
		}
	})

	t.Run("bool value false", func(t *testing.T) {
		result := Bool(false)
		expected := sql.NullBool{Bool: false, Valid: true}
		if result != expected {
			t.Errorf("Bool(false) = %v, want %v", result, expected)
		}
	})

	t.Run("bool pointer true", func(t *testing.T) {
		b := true
		result := Bool(&b)
		expected := sql.NullBool{Bool: true, Valid: true}
		if result != expected {
			t.Errorf("Bool(&true) = %v, want %v", result, expected)
		}
	})

	t.Run("bool pointer false", func(t *testing.T) {
		b := false
		result := Bool(&b)
		expected := sql.NullBool{Bool: false, Valid: true}
		if result != expected {
			t.Errorf("Bool(&false) = %v, want %v", result, expected)
		}
	})

	t.Run("nil bool pointer", func(t *testing.T) {
		result := Bool((*bool)(nil))
		expected := sql.NullBool{Bool: false, Valid: false}
		if result != expected {
			t.Errorf("Bool((*bool)(nil)) = %v, want %v", result, expected)
		}
	})
}

func TestBoolGenericTypeConstraints(t *testing.T) {
	// Test that the generic function works with both bool and *bool types
	t.Run("generic with bool", func(t *testing.T) {
		var b bool = true
		result := Bool(b)
		expected := sql.NullBool{Bool: true, Valid: true}
		if result != expected {
			t.Errorf("Bool[bool](true) = %v, want %v", result, expected)
		}
	})

	t.Run("generic with *bool", func(t *testing.T) {
		b := false
		var ptr *bool = &b
		result := Bool(ptr)
		expected := sql.NullBool{Bool: false, Valid: true}
		if result != expected {
			t.Errorf("Bool[*bool](&false) = %v, want %v", result, expected)
		}
	})
}

// Benchmark tests
func BenchmarkBoolAny(b *testing.B) {
	nullBool := sql.NullBool{Bool: true, Valid: true}
	for i := 0; i < b.N; i++ {
		BoolAny(nullBool)
	}
}

func BenchmarkBoolNil(b *testing.B) {
	nullBool := sql.NullBool{Bool: true, Valid: true}
	for i := 0; i < b.N; i++ {
		BoolNil(nullBool)
	}
}

func BenchmarkBool(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Bool(true)
	}
}
