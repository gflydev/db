package null

import (
	"database/sql"
	"testing"
)

func TestByteNil(t *testing.T) {
	tests := []struct {
		name     string
		input    sql.NullByte
		expected *byte
	}{
		{
			name:     "valid byte 65 (A)",
			input:    sql.NullByte{Byte: 65, Valid: true},
			expected: func() *byte { b := byte(65); return &b }(),
		},
		{
			name:     "valid byte 0",
			input:    sql.NullByte{Byte: 0, Valid: true},
			expected: func() *byte { b := byte(0); return &b }(),
		},
		{
			name:     "valid byte 255",
			input:    sql.NullByte{Byte: 255, Valid: true},
			expected: func() *byte { b := byte(255); return &b }(),
		},
		{
			name:     "invalid null",
			input:    sql.NullByte{Byte: 0, Valid: false},
			expected: nil,
		},
		{
			name:     "invalid null with value",
			input:    sql.NullByte{Byte: 100, Valid: false},
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ByteNil(tt.input)
			if tt.expected == nil {
				if result != nil {
					t.Errorf("ByteNil(%v) = %v, want nil", tt.input, result)
				}
			} else {
				if result == nil {
					t.Errorf("ByteNil(%v) = nil, want %v", tt.input, *tt.expected)
				} else if *result != *tt.expected {
					t.Errorf("ByteNil(%v) = %v, want %v", tt.input, *result, *tt.expected)
				}
			}
		})
	}
}

func TestByteVal(t *testing.T) {
	tests := []struct {
		name     string
		input    sql.NullByte
		expected byte
	}{
		{
			name:     "valid byte 65 (A)",
			input:    sql.NullByte{Byte: 65, Valid: true},
			expected: byte(65),
		},
		{
			name:     "valid byte 0",
			input:    sql.NullByte{Byte: 0, Valid: true},
			expected: byte(0),
		},
		{
			name:     "valid byte 255",
			input:    sql.NullByte{Byte: 255, Valid: true},
			expected: byte(255),
		},
		{
			name:     "invalid null",
			input:    sql.NullByte{Byte: 0, Valid: false},
			expected: byte(0),
		},
		{
			name:     "invalid null with value",
			input:    sql.NullByte{Byte: 100, Valid: false},
			expected: byte(0),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ByteVal(tt.input)
			if result != tt.expected {
				t.Errorf("ByteVal(%v) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestByte(t *testing.T) {
	t.Run("byte value 65 (A)", func(t *testing.T) {
		result := Byte(byte(65))
		expected := sql.NullByte{Byte: 65, Valid: true}
		if result != expected {
			t.Errorf("Byte(65) = %v, want %v", result, expected)
		}
	})

	t.Run("byte value 0", func(t *testing.T) {
		result := Byte(byte(0))
		expected := sql.NullByte{Byte: 0, Valid: true}
		if result != expected {
			t.Errorf("Byte(0) = %v, want %v", result, expected)
		}
	})

	t.Run("byte value 255", func(t *testing.T) {
		result := Byte(byte(255))
		expected := sql.NullByte{Byte: 255, Valid: true}
		if result != expected {
			t.Errorf("Byte(255) = %v, want %v", result, expected)
		}
	})

	t.Run("byte pointer 97 (a)", func(t *testing.T) {
		b := byte(97)
		result := Byte(&b)
		expected := sql.NullByte{Byte: 97, Valid: true}
		if result != expected {
			t.Errorf("Byte(&97) = %v, want %v", result, expected)
		}
	})

	t.Run("byte pointer 0", func(t *testing.T) {
		b := byte(0)
		result := Byte(&b)
		expected := sql.NullByte{Byte: 0, Valid: true}
		if result != expected {
			t.Errorf("Byte(&0) = %v, want %v", result, expected)
		}
	})

	t.Run("nil byte pointer", func(t *testing.T) {
		result := Byte((*byte)(nil))
		expected := sql.NullByte{Byte: 0, Valid: false}
		if result != expected {
			t.Errorf("Byte((*byte)(nil)) = %v, want %v", result, expected)
		}
	})
}

func TestByteGenericTypeConstraints(t *testing.T) {
	// Test that the generic function works with both byte and *byte types
	t.Run("generic with byte", func(t *testing.T) {
		var b byte = 42
		result := Byte(b)
		expected := sql.NullByte{Byte: 42, Valid: true}
		if result != expected {
			t.Errorf("Byte[byte](42) = %v, want %v", result, expected)
		}
	})

	t.Run("generic with *byte", func(t *testing.T) {
		b := byte(123)
		var ptr *byte = &b
		result := Byte(ptr)
		expected := sql.NullByte{Byte: 123, Valid: true}
		if result != expected {
			t.Errorf("Byte[*byte](&123) = %v, want %v", result, expected)
		}
	})
}

func TestByteEdgeCases(t *testing.T) {
	t.Run("minimum byte value", func(t *testing.T) {
		result := Byte(byte(0))
		expected := sql.NullByte{Byte: 0, Valid: true}
		if result != expected {
			t.Errorf("Byte(0) = %v, want %v", result, expected)
		}
	})

	t.Run("maximum byte value", func(t *testing.T) {
		result := Byte(byte(255))
		expected := sql.NullByte{Byte: 255, Valid: true}
		if result != expected {
			t.Errorf("Byte(255) = %v, want %v", result, expected)
		}
	})

	t.Run("common ASCII values", func(t *testing.T) {
		testCases := []struct {
			value    byte
			expected sql.NullByte
		}{
			{32, sql.NullByte{Byte: 32, Valid: true}},   // space
			{48, sql.NullByte{Byte: 48, Valid: true}},   // '0'
			{65, sql.NullByte{Byte: 65, Valid: true}},   // 'A'
			{97, sql.NullByte{Byte: 97, Valid: true}},   // 'a'
			{126, sql.NullByte{Byte: 126, Valid: true}}, // '~'
		}

		for _, tc := range testCases {
			result := Byte(tc.value)
			if result != tc.expected {
				t.Errorf("Byte(%d) = %v, want %v", tc.value, result, tc.expected)
			}
		}
	})
}

// Benchmark tests

func BenchmarkByteNil(b *testing.B) {
	nullByte := sql.NullByte{Byte: 65, Valid: true}
	for i := 0; i < b.N; i++ {
		ByteNil(nullByte)
	}
}

func BenchmarkByte(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Byte(byte(65))
	}
}

func BenchmarkBytePointer(b *testing.B) {
	byteVal := byte(65)
	for i := 0; i < b.N; i++ {
		Byte(&byteVal)
	}
}
