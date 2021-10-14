package util

import (
	"time"
)

// ConverterFunc is a function that maps from a value from one source type to destination type
// ConverterFunc should always return a valid dst value
type ConverterFunc func(src interface{}) (dst interface{})

func ConvertRFC3339MilliToTime(src interface{}) interface{} {
	RFC3339Milli := "2006-01-02T15:04:05.000Z07:00"
	dst, err := time.Parse(RFC3339Milli, src.(string))

	if err != nil {
		return time.Time{}
	}

	return dst
}

func ConvertTimeToRFC3339Milli(src interface{}) interface{} {
	RFC3339Milli := "2006-01-02T15:04:05.000Z07:00"
	dst := src.(time.Time).Format(RFC3339Milli)
	return dst
}
