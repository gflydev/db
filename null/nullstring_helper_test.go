package null

import (
	"database/sql"
	"strings"
	"testing"
)

func TestStringNil(t *testing.T) {
	tests := []struct {
		name     string
		input    sql.NullString
		expected *string
	}{
		{
			name:     "valid non-empty string",
			input:    sql.NullString{String: "hello", Valid: true},
			expected: func() *string { s := "hello"; return &s }(),
		},
		{
			name:     "valid empty string",
			input:    sql.NullString{String: "", Valid: true},
			expected: func() *string { s := ""; return &s }(),
		},
		{
			name:     "valid string with spaces",
			input:    sql.NullString{String: "  hello world  ", Valid: true},
			expected: func() *string { s := "  hello world  "; return &s }(),
		},
		{
			name:     "valid string with special characters",
			input:    sql.NullString{String: "hello\nworld\t!", Valid: true},
			expected: func() *string { s := "hello\nworld\t!"; return &s }(),
		},
		{
			name:     "valid unicode string",
			input:    sql.NullString{String: "こんにちは世界", Valid: true},
			expected: func() *string { s := "こんにちは世界"; return &s }(),
		},
		{
			name:     "invalid null",
			input:    sql.NullString{String: "", Valid: false},
			expected: nil,
		},
		{
			name:     "invalid null with value",
			input:    sql.NullString{String: "hello", Valid: false},
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := StringNil(tt.input)
			if tt.expected == nil {
				if result != nil {
					t.Errorf("StringNil(%v) = %v, want nil", tt.input, result)
				}
			} else {
				if result == nil {
					t.Errorf("StringNil(%v) = nil, want %v", tt.input, *tt.expected)
				} else if *result != *tt.expected {
					t.Errorf("StringNil(%v) = %v, want %v", tt.input, *result, *tt.expected)
				}
			}
		})
	}
}

func TestStringVal(t *testing.T) {
	tests := []struct {
		name     string
		input    sql.NullString
		expected string
	}{
		{
			name:     "valid string hello",
			input:    sql.NullString{String: "hello", Valid: true},
			expected: "hello",
		},
		{
			name:     "valid empty string",
			input:    sql.NullString{String: "", Valid: true},
			expected: "",
		},
		{
			name:     "valid string with spaces",
			input:    sql.NullString{String: "  hello world  ", Valid: true},
			expected: "  hello world  ",
		},
		{
			name:     "valid string with special characters",
			input:    sql.NullString{String: "hello\nworld\t!", Valid: true},
			expected: "hello\nworld\t!",
		},
		{
			name:     "invalid null",
			input:    sql.NullString{String: "", Valid: false},
			expected: "",
		},
		{
			name:     "invalid null with value",
			input:    sql.NullString{String: "test", Valid: false},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := StringVal(tt.input)
			if result != tt.expected {
				t.Errorf("StringVal(%v) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestString(t *testing.T) {
	t.Run("string value non-empty", func(t *testing.T) {
		result := String("hello")
		expected := sql.NullString{String: "hello", Valid: true}
		if result != expected {
			t.Errorf("String(\"hello\") = %v, want %v", result, expected)
		}
	})

	t.Run("string value empty", func(t *testing.T) {
		result := String("")
		expected := sql.NullString{String: "", Valid: true}
		if result != expected {
			t.Errorf("String(\"\") = %v, want %v", result, expected)
		}
	})

	t.Run("string value with spaces", func(t *testing.T) {
		result := String("  hello world  ")
		expected := sql.NullString{String: "  hello world  ", Valid: true}
		if result != expected {
			t.Errorf("String(\"  hello world  \") = %v, want %v", result, expected)
		}
	})

	t.Run("string pointer non-empty", func(t *testing.T) {
		s := "world"
		result := String(&s)
		expected := sql.NullString{String: "world", Valid: true}
		if result != expected {
			t.Errorf("String(&\"world\") = %v, want %v", result, expected)
		}
	})

	t.Run("string pointer empty", func(t *testing.T) {
		s := ""
		result := String(&s)
		expected := sql.NullString{String: "", Valid: true}
		if result != expected {
			t.Errorf("String(&\"\") = %v, want %v", result, expected)
		}
	})

	t.Run("nil string pointer", func(t *testing.T) {
		result := String((*string)(nil))
		expected := sql.NullString{String: "", Valid: false}
		if result != expected {
			t.Errorf("String((*string)(nil)) = %v, want %v", result, expected)
		}
	})
}

func TestStringSpecialCases(t *testing.T) {
	t.Run("string with newlines", func(t *testing.T) {
		input := "line1\nline2\nline3"
		result := String(input)
		expected := sql.NullString{String: input, Valid: true}
		if result != expected {
			t.Errorf("String with newlines = %v, want %v", result, expected)
		}
	})

	t.Run("string with tabs", func(t *testing.T) {
		input := "col1\tcol2\tcol3"
		result := String(input)
		expected := sql.NullString{String: input, Valid: true}
		if result != expected {
			t.Errorf("String with tabs = %v, want %v", result, expected)
		}
	})

	t.Run("string with unicode", func(t *testing.T) {
		input := "Hello 世界 🌍"
		result := String(input)
		expected := sql.NullString{String: input, Valid: true}
		if result != expected {
			t.Errorf("String with unicode = %v, want %v", result, expected)
		}
	})

	t.Run("very long string", func(t *testing.T) {
		input := strings.Repeat("abcdefghij", 100) // 1000 characters
		result := String(input)
		expected := sql.NullString{String: input, Valid: true}
		if result != expected {
			t.Errorf("Very long string test failed")
		}
	})

	t.Run("string with null character", func(t *testing.T) {
		input := "hello\x00world"
		result := String(input)
		expected := sql.NullString{String: input, Valid: true}
		if result != expected {
			t.Errorf("String with null character = %v, want %v", result, expected)
		}
	})

	t.Run("string with only whitespace", func(t *testing.T) {
		input := "   \t\n\r   "
		result := String(input)
		expected := sql.NullString{String: input, Valid: true}
		if result != expected {
			t.Errorf("String with only whitespace = %v, want %v", result, expected)
		}
	})
}

func TestStringGenericTypeConstraints(t *testing.T) {
	// Test that the generic function works with both string and *string types
	t.Run("generic with string", func(t *testing.T) {
		var s string = "test"
		result := String(s)
		expected := sql.NullString{String: "test", Valid: true}
		if result != expected {
			t.Errorf("String[string](\"test\") = %v, want %v", result, expected)
		}
	})

	t.Run("generic with *string", func(t *testing.T) {
		s := "pointer test"
		var ptr *string = &s
		result := String(ptr)
		expected := sql.NullString{String: "pointer test", Valid: true}
		if result != expected {
			t.Errorf("String[*string](&\"pointer test\") = %v, want %v", result, expected)
		}
	})
}

func TestStringEdgeCases(t *testing.T) {
	t.Run("single character", func(t *testing.T) {
		result := String("a")
		expected := sql.NullString{String: "a", Valid: true}
		if result != expected {
			t.Errorf("String(\"a\") = %v, want %v", result, expected)
		}
	})

	t.Run("numeric string", func(t *testing.T) {
		result := String("12345")
		expected := sql.NullString{String: "12345", Valid: true}
		if result != expected {
			t.Errorf("String(\"12345\") = %v, want %v", result, expected)
		}
	})

	t.Run("boolean string", func(t *testing.T) {
		result := String("true")
		expected := sql.NullString{String: "true", Valid: true}
		if result != expected {
			t.Errorf("String(\"true\") = %v, want %v", result, expected)
		}
	})

	t.Run("json string", func(t *testing.T) {
		input := `{"key": "value", "number": 42}`
		result := String(input)
		expected := sql.NullString{String: input, Valid: true}
		if result != expected {
			t.Errorf("String with JSON = %v, want %v", result, expected)
		}
	})

	t.Run("sql injection attempt", func(t *testing.T) {
		input := "'; DROP TABLE users; --"
		result := String(input)
		expected := sql.NullString{String: input, Valid: true}
		if result != expected {
			t.Errorf("String with SQL injection = %v, want %v", result, expected)
		}
	})
}

// Benchmark tests

func BenchmarkStringNil(b *testing.B) {
	nullString := sql.NullString{String: "hello world", Valid: true}
	for i := 0; i < b.N; i++ {
		StringNil(nullString)
	}
}

func BenchmarkString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		String("hello world")
	}
}

func BenchmarkStringPointer(b *testing.B) {
	str := "hello world"
	for i := 0; i < b.N; i++ {
		String(&str)
	}
}

func BenchmarkStringLong(b *testing.B) {
	longString := strings.Repeat("abcdefghij", 100)
	for i := 0; i < b.N; i++ {
		String(longString)
	}
}
