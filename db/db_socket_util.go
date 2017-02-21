package db

import (
	"database/sql"
	"errors"
	"github.com/lib/pq"
	"time"
)

func GetStringValue(value sql.NullString, default_value string) string {
	if value.Valid {
		return value.String
	} else {
		return default_value
	}
}

func GetInt64Value(value sql.NullInt64, default_value int64) int64 {
	if value.Valid {
		return value.Int64
	} else {
		return default_value
	}
}

func GetFloat64Value(value sql.NullFloat64, default_value float64) float64 {
	if value.Valid {
		return value.Float64
	} else {
		return default_value
	}
}

func GetTimeValue(value pq.NullTime) (time.Time, error) {
	if value.Valid {
		return value.Time, nil
	}

	return time.Now(), errors.New("Time is not valid")

}
