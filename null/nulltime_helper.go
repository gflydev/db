package null

import (
	"database/sql"
	"database/sql/driver"
	"time"
)

// ScanTime function will scan NullTime value.
func ScanTime(nullTime sql.NullTime) driver.Value {
	if !nullTime.Valid {
		return nil
	}
	return nullTime.Time
}

// Time function will create a NullTime object.
func Time(val time.Time) sql.NullTime {
	return sql.NullTime{
		Time:  val,
		Valid: true,
	}
}

// NowTime function will create a NullTime object.
func NowTime() sql.NullTime {
	return Time(time.Now())
}

// TimeVal function will scan NullTime value.
func TimeVal(nullTime sql.NullTime) *time.Time {
	if !nullTime.Valid {
		return nil
	}
	return &nullTime.Time
}
