package null

import (
	"database/sql"
	"database/sql/driver"
	"time"
)

// TimeAny function will scan NullTime value.
func TimeAny(nullTime sql.NullTime) driver.Value {
	if !nullTime.Valid {
		return nil
	}
	return nullTime.Time
}

// TimeNil function will scan NullTime value.
func TimeNil(nullTime sql.NullTime) *time.Time {
	if !nullTime.Valid {
		return nil
	}
	return &nullTime.Time
}

// Time function will create a NullTime object.
// It accepts both time.Time and *time.Time values.
func Time(val any) sql.NullTime {
	switch v := val.(type) {
	case time.Time:
		return sql.NullTime{
			Time:  v,
			Valid: true,
		}
	case *time.Time:
		if v == nil {
			return sql.NullTime{
				Time:  time.Time{},
				Valid: false,
			}
		}
		return sql.NullTime{
			Time:  *v,
			Valid: true,
		}
	default:
		// For any other type, return invalid NullTime
		return sql.NullTime{
			Time:  time.Time{},
			Valid: false,
		}
	}
}

// NowTime function will create a NullTime object.
func NowTime() sql.NullTime {
	return Time(time.Now())
}
