package null

import (
	"database/sql"
	"database/sql/driver"
	"testing"
	"time"
)

func TestTimeAny(t *testing.T) {
	now := time.Now()
	utc := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)
	zeroTime := time.Time{}

	tests := []struct {
		name     string
		input    sql.NullTime
		expected driver.Value
	}{
		{
			name:     "valid current time",
			input:    sql.NullTime{Time: now, Valid: true},
			expected: now,
		},
		{
			name:     "valid UTC time",
			input:    sql.NullTime{Time: utc, Valid: true},
			expected: utc,
		},
		{
			name:     "valid zero time",
			input:    sql.NullTime{Time: zeroTime, Valid: true},
			expected: zeroTime,
		},
		{
			name:     "valid unix epoch",
			input:    sql.NullTime{Time: time.Unix(0, 0), Valid: true},
			expected: time.Unix(0, 0),
		},
		{
			name:     "valid far future",
			input:    sql.NullTime{Time: time.Date(2099, 12, 31, 23, 59, 59, 999999999, time.UTC), Valid: true},
			expected: time.Date(2099, 12, 31, 23, 59, 59, 999999999, time.UTC),
		},
		{
			name:     "valid far past",
			input:    sql.NullTime{Time: time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC), Valid: true},
			expected: time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "invalid null",
			input:    sql.NullTime{Time: zeroTime, Valid: false},
			expected: nil,
		},
		{
			name:     "invalid null with time value",
			input:    sql.NullTime{Time: now, Valid: false},
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := TimeAny(tt.input)
			if result != tt.expected {
				t.Errorf("TimeAny(%v) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestTimeNil(t *testing.T) {
	now := time.Now()
	utc := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)
	zeroTime := time.Time{}

	tests := []struct {
		name     string
		input    sql.NullTime
		expected *time.Time
	}{
		{
			name:     "valid current time",
			input:    sql.NullTime{Time: now, Valid: true},
			expected: &now,
		},
		{
			name:     "valid UTC time",
			input:    sql.NullTime{Time: utc, Valid: true},
			expected: &utc,
		},
		{
			name:     "valid zero time",
			input:    sql.NullTime{Time: zeroTime, Valid: true},
			expected: &zeroTime,
		},
		{
			name:     "valid unix epoch",
			input:    sql.NullTime{Time: time.Unix(0, 0), Valid: true},
			expected: func() *time.Time { t := time.Unix(0, 0); return &t }(),
		},
		{
			name:     "invalid null",
			input:    sql.NullTime{Time: zeroTime, Valid: false},
			expected: nil,
		},
		{
			name:     "invalid null with time value",
			input:    sql.NullTime{Time: now, Valid: false},
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := TimeNil(tt.input)
			if tt.expected == nil {
				if result != nil {
					t.Errorf("TimeNil(%v) = %v, want nil", tt.input, result)
				}
			} else {
				if result == nil {
					t.Errorf("TimeNil(%v) = nil, want %v", tt.input, *tt.expected)
				} else if !result.Equal(*tt.expected) {
					t.Errorf("TimeNil(%v) = %v, want %v", tt.input, *result, *tt.expected)
				}
			}
		})
	}
}

func TestTime(t *testing.T) {
	now := time.Now()
	utc := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)
	zeroTime := time.Time{}

	t.Run("time.Time value current", func(t *testing.T) {
		result := Time(now)
		expected := sql.NullTime{Time: now, Valid: true}
		if !result.Time.Equal(expected.Time) || result.Valid != expected.Valid {
			t.Errorf("Time(now) = %v, want %v", result, expected)
		}
	})

	t.Run("time.Time value UTC", func(t *testing.T) {
		result := Time(utc)
		expected := sql.NullTime{Time: utc, Valid: true}
		if !result.Time.Equal(expected.Time) || result.Valid != expected.Valid {
			t.Errorf("Time(utc) = %v, want %v", result, expected)
		}
	})

	t.Run("time.Time value zero", func(t *testing.T) {
		result := Time(zeroTime)
		expected := sql.NullTime{Time: zeroTime, Valid: true}
		if !result.Time.Equal(expected.Time) || result.Valid != expected.Valid {
			t.Errorf("Time(zeroTime) = %v, want %v", result, expected)
		}
	})

	t.Run("time.Time pointer current", func(t *testing.T) {
		result := Time(&now)
		expected := sql.NullTime{Time: now, Valid: true}
		if !result.Time.Equal(expected.Time) || result.Valid != expected.Valid {
			t.Errorf("Time(&now) = %v, want %v", result, expected)
		}
	})

	t.Run("time.Time pointer UTC", func(t *testing.T) {
		result := Time(&utc)
		expected := sql.NullTime{Time: utc, Valid: true}
		if !result.Time.Equal(expected.Time) || result.Valid != expected.Valid {
			t.Errorf("Time(&utc) = %v, want %v", result, expected)
		}
	})

	t.Run("time.Time pointer zero", func(t *testing.T) {
		result := Time(&zeroTime)
		expected := sql.NullTime{Time: zeroTime, Valid: true}
		if !result.Time.Equal(expected.Time) || result.Valid != expected.Valid {
			t.Errorf("Time(&zeroTime) = %v, want %v", result, expected)
		}
	})

	t.Run("nil time.Time pointer", func(t *testing.T) {
		result := Time((*time.Time)(nil))
		expected := sql.NullTime{Time: time.Time{}, Valid: false}
		if !result.Time.Equal(expected.Time) || result.Valid != expected.Valid {
			t.Errorf("Time((*time.Time)(nil)) = %v, want %v", result, expected)
		}
	})
}

func TestTimeSpecialCases(t *testing.T) {
	t.Run("unix epoch", func(t *testing.T) {
		epoch := time.Unix(0, 0)
		result := Time(epoch)
		expected := sql.NullTime{Time: epoch, Valid: true}
		if !result.Time.Equal(expected.Time) || result.Valid != expected.Valid {
			t.Errorf("Time(epoch) = %v, want %v", result, expected)
		}
	})

	t.Run("far future date", func(t *testing.T) {
		future := time.Date(2099, 12, 31, 23, 59, 59, 999999999, time.UTC)
		result := Time(future)
		expected := sql.NullTime{Time: future, Valid: true}
		if !result.Time.Equal(expected.Time) || result.Valid != expected.Valid {
			t.Errorf("Time(future) = %v, want %v", result, expected)
		}
	})

	t.Run("far past date", func(t *testing.T) {
		past := time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC)
		result := Time(past)
		expected := sql.NullTime{Time: past, Valid: true}
		if !result.Time.Equal(expected.Time) || result.Valid != expected.Valid {
			t.Errorf("Time(past) = %v, want %v", result, expected)
		}
	})

	t.Run("different timezone", func(t *testing.T) {
		loc, _ := time.LoadLocation("America/New_York")
		nyTime := time.Date(2023, 6, 15, 14, 30, 0, 0, loc)
		result := Time(nyTime)
		expected := sql.NullTime{Time: nyTime, Valid: true}
		if !result.Time.Equal(expected.Time) || result.Valid != expected.Valid {
			t.Errorf("Time(nyTime) = %v, want %v", result, expected)
		}
	})

	t.Run("leap year date", func(t *testing.T) {
		leapDay := time.Date(2020, 2, 29, 12, 0, 0, 0, time.UTC)
		result := Time(leapDay)
		expected := sql.NullTime{Time: leapDay, Valid: true}
		if !result.Time.Equal(expected.Time) || result.Valid != expected.Valid {
			t.Errorf("Time(leapDay) = %v, want %v", result, expected)
		}
	})

	t.Run("nanosecond precision", func(t *testing.T) {
		nanoTime := time.Date(2023, 1, 1, 12, 0, 0, 123456789, time.UTC)
		result := Time(nanoTime)
		expected := sql.NullTime{Time: nanoTime, Valid: true}
		if !result.Time.Equal(expected.Time) || result.Valid != expected.Valid {
			t.Errorf("Time(nanoTime) = %v, want %v", result, expected)
		}
	})
}

func TestTimeNow(t *testing.T) {
	t.Run("TimeNow creates valid NullTime", func(t *testing.T) {
		before := time.Now()
		result := TimeNow()
		after := time.Now()

		if !result.Valid {
			t.Errorf("TimeNow() should create valid NullTime, got Valid = %v", result.Valid)
		}

		if result.Time.Before(before) || result.Time.After(after) {
			t.Errorf("TimeNow() time %v should be between %v and %v", result.Time, before, after)
		}
	})

	t.Run("TimeNow equivalent to Time(time.Now())", func(t *testing.T) {
		// We can't test exact equality due to timing, but we can test structure
		result1 := TimeNow()
		now := time.Now()
		result2 := Time(now)

		if result1.Valid != result2.Valid {
			t.Errorf("TimeNow().Valid = %v, Time(time.Now()).Valid = %v", result1.Valid, result2.Valid)
		}

		// Both should be valid
		if !result1.Valid || !result2.Valid {
			t.Errorf("Both TimeNow() and Time(time.Now()) should be valid")
		}
	})

	t.Run("multiple TimeNow calls", func(t *testing.T) {
		result1 := TimeNow()
		time.Sleep(1 * time.Millisecond) // Small delay to ensure different times
		result2 := TimeNow()

		if !result1.Valid || !result2.Valid {
			t.Errorf("Both TimeNow() calls should be valid")
		}

		if !result2.Time.After(result1.Time) && !result2.Time.Equal(result1.Time) {
			t.Errorf("Second TimeNow() call should be after or equal to first")
		}
	})
}

func TestTimeGenericTypeConstraints(t *testing.T) {
	now := time.Now()

	// Test that the generic function works with both time.Time and *time.Time types
	t.Run("generic with time.Time", func(t *testing.T) {
		var tm time.Time = now
		result := Time(tm)
		expected := sql.NullTime{Time: now, Valid: true}
		if !result.Time.Equal(expected.Time) || result.Valid != expected.Valid {
			t.Errorf("Time[time.Time](now) = %v, want %v", result, expected)
		}
	})

	t.Run("generic with *time.Time", func(t *testing.T) {
		var ptr *time.Time = &now
		result := Time(ptr)
		expected := sql.NullTime{Time: now, Valid: true}
		if !result.Time.Equal(expected.Time) || result.Valid != expected.Valid {
			t.Errorf("Time[*time.Time](&now) = %v, want %v", result, expected)
		}
	})
}

func TestTimeEdgeCases(t *testing.T) {
	t.Run("time with different locations", func(t *testing.T) {
		utc := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)
		local := utc.Local()

		resultUTC := Time(utc)
		resultLocal := Time(local)

		if !resultUTC.Valid || !resultLocal.Valid {
			t.Errorf("Both UTC and local times should be valid")
		}

		// They represent the same instant but in different locations
		if !resultUTC.Time.Equal(resultLocal.Time) {
			t.Errorf("UTC and local times should represent the same instant")
		}
	})

	t.Run("time parsing from string", func(t *testing.T) {
		timeStr := "2023-01-01T12:00:00Z"
		parsed, err := time.Parse(time.RFC3339, timeStr)
		if err != nil {
			t.Fatalf("Failed to parse time: %v", err)
		}

		result := Time(parsed)
		expected := sql.NullTime{Time: parsed, Valid: true}
		if !result.Time.Equal(expected.Time) || result.Valid != expected.Valid {
			t.Errorf("Time(parsed) = %v, want %v", result, expected)
		}
	})

	t.Run("time with microsecond precision", func(t *testing.T) {
		microTime := time.Date(2023, 1, 1, 12, 0, 0, 123456000, time.UTC)
		result := Time(microTime)
		expected := sql.NullTime{Time: microTime, Valid: true}
		if !result.Time.Equal(expected.Time) || result.Valid != expected.Valid {
			t.Errorf("Time(microTime) = %v, want %v", result, expected)
		}
	})
}

// Benchmark tests
func BenchmarkTimeAny(b *testing.B) {
	nullTime := sql.NullTime{Time: time.Now(), Valid: true}
	for i := 0; i < b.N; i++ {
		TimeAny(nullTime)
	}
}

func BenchmarkTimeNil(b *testing.B) {
	nullTime := sql.NullTime{Time: time.Now(), Valid: true}
	for i := 0; i < b.N; i++ {
		TimeNil(nullTime)
	}
}

func BenchmarkTime(b *testing.B) {
	now := time.Now()
	for i := 0; i < b.N; i++ {
		Time(now)
	}
}

func BenchmarkTimePointer(b *testing.B) {
	now := time.Now()
	for i := 0; i < b.N; i++ {
		Time(&now)
	}
}

func BenchmarkTimeNow(b *testing.B) {
	for i := 0; i < b.N; i++ {
		TimeNow()
	}
}

func BenchmarkTimeNowVsTimeTimeNow(b *testing.B) {
	b.Run("TimeNow", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			TimeNow()
		}
	})

	b.Run("Time(time.Now())", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			Time(time.Now())
		}
	})
}
