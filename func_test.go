package db

import (
	"encoding/json"
	"testing"
)

// TestToStr verifies string conversion across the supported scalar types,
// including the integer widths that were previously unhandled.
func TestToStr(t *testing.T) {
	tests := []struct {
		name     string
		input    any
		expected string
	}{
		{"bool true", true, "true"},
		{"int", 42, "42"},
		{"int8", int8(8), "8"},
		{"int16", int16(16), "16"},
		{"int32", int32(32), "32"},
		{"int64", int64(64), "64"},
		{"uint", uint(42), "42"},
		{"uint8", uint8(8), "8"},
		{"uint16", uint16(16), "16"},
		{"uint32", uint32(32), "32"},
		{"uint64", uint64(64), "64"},
		{"float64", 3.14, "3.140000"},
		{"string", "hello", "hello"},
		{"bytes", []byte("world"), "world"},
		{"json number", json.Number("123.45"), "123.45"},
		{"unsupported returns empty", struct{}{}, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := toStr(tt.input); got != tt.expected {
				t.Errorf("toStr(%v) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

// setValueTarget exercises setValue across several field kinds.
type setValueTarget struct {
	Name   string
	Age    int
	Small  int32
	Tiny   uint8
	Active bool
	Score  float64
}

// TestSetValue verifies dynamic field assignment with type conversion,
// covering narrow integer types that depend on toStr working correctly.
func TestSetValue(t *testing.T) {
	target := &setValueTarget{}

	if err := setValue(target, "Name", "Vinh"); err != nil {
		t.Fatalf("Name: %v", err)
	}
	if err := setValue(target, "Age", int64(30)); err != nil {
		t.Fatalf("Age: %v", err)
	}
	// int32 field fed a narrow int value: previously toStr(int32) returned ""
	// and the field silently became 0.
	if err := setValue(target, "Small", int32(7)); err != nil {
		t.Fatalf("Small: %v", err)
	}
	if err := setValue(target, "Tiny", uint8(9)); err != nil {
		t.Fatalf("Tiny: %v", err)
	}
	if err := setValue(target, "Active", true); err != nil {
		t.Fatalf("Active: %v", err)
	}
	if err := setValue(target, "Score", 4.5); err != nil {
		t.Fatalf("Score: %v", err)
	}

	if target.Name != "Vinh" {
		t.Errorf("Name = %q, want Vinh", target.Name)
	}
	if target.Age != 30 {
		t.Errorf("Age = %d, want 30", target.Age)
	}
	if target.Small != 7 {
		t.Errorf("Small = %d, want 7", target.Small)
	}
	if target.Tiny != 9 {
		t.Errorf("Tiny = %d, want 9", target.Tiny)
	}
	if !target.Active {
		t.Errorf("Active = false, want true")
	}
	if target.Score != 4.5 {
		t.Errorf("Score = %v, want 4.5", target.Score)
	}
}

// TestSetValueUnknownType verifies an unsupported field kind returns an error.
func TestSetValueUnknownType(t *testing.T) {
	target := &struct{ Data []string }{}
	if err := setValue(target, "Data", []string{"x"}); err == nil {
		t.Error("expected error for unsupported field type")
	}
}
