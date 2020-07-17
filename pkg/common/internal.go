package common

import (
	"errors"
	"math"
	"strconv"
	"strings"
	"time"
)

// ParseTimestamp parses a Time according to the standard Telegraf options.
// These are generally displayed in the toml similar to:
//   json_time_key= "timestamp"
//   json_time_format = "2006-01-02T15:04:05Z07:00"
//   json_timezone = "America/Los_Angeles"
//
// The format can be one of "unix", "unix_ms", "unix_us", "unix_ns", or a Go
// time layout suitable for time.Parse.
//
// When using the "unix" format, a optional fractional component is allowed.
// Specific unix time precisions cannot have a fractional component.
//
// Unix times may be an int64, float64, or string.  When using a Go format
// string the timestamp must be a string.
//
// The location is a location string suitable for time.LoadLocation.  Unix
// times do not use the location string, a unix time is always return in the
// UTC location.
func ParseTimestamp(format string, timestamp interface{}, location string) (time.Time, error) {
	switch format {
	case "unix", "unix_ms", "unix_us", "unix_ns":
		return parseUnix(format, timestamp)
	default:
		if location == "" {
			location = "UTC"
		}
		return parseTime(format, timestamp, location)
	}
}

func parseUnix(format string, timestamp interface{}) (time.Time, error) {
	integer, fractional, err := parseComponents(timestamp)
	if err != nil {
		return time.Unix(0, 0), err
	}

	switch strings.ToLower(format) {
	case "unix":
		return time.Unix(integer, fractional).UTC(), nil
	case "unix_ms":
		return time.Unix(0, integer*1e6).UTC(), nil
	case "unix_us":
		return time.Unix(0, integer*1e3).UTC(), nil
	case "unix_ns":
		return time.Unix(0, integer).UTC(), nil
	default:
		return time.Unix(0, 0), errors.New("unsupported type")
	}
}

// Returns the integers before and after an optional decimal point.  Both '.'
// and ',' are supported for the decimal point.  The timestamp can be an int64,
// float64, or string.
//   ex: "42.5" -> (42, 5, nil)
func parseComponents(timestamp interface{}) (int64, int64, error) {
	switch ts := timestamp.(type) {
	case string:
		parts := strings.SplitN(ts, ".", 2)
		if len(parts) == 2 {
			return parseUnixTimeComponents(parts[0], parts[1])
		}

		parts = strings.SplitN(ts, ",", 2)
		if len(parts) == 2 {
			return parseUnixTimeComponents(parts[0], parts[1])
		}

		integer, err := strconv.ParseInt(ts, 10, 64)
		if err != nil {
			return 0, 0, err
		}
		return integer, 0, nil
	case int64:
		return ts, 0, nil
	case float64:
		integer, fractional := math.Modf(ts)
		return int64(integer), int64(fractional * 1e9), nil
	default:
		return 0, 0, errors.New("unsupported type")
	}
}

func parseUnixTimeComponents(first, second string) (int64, int64, error) {
	integer, err := strconv.ParseInt(first, 10, 64)
	if err != nil {
		return 0, 0, err
	}

	// Convert to nanoseconds, dropping any greater precision.
	buf := []byte("000000000")
	copy(buf, second)

	fractional, err := strconv.ParseInt(string(buf), 10, 64)
	if err != nil {
		return 0, 0, err
	}
	return integer, fractional, nil
}

// ParseTime parses a string timestamp according to the format string.
func parseTime(format string, timestamp interface{}, location string) (time.Time, error) {
	switch ts := timestamp.(type) {
	case string:
		loc, err := time.LoadLocation(location)
		if err != nil {
			return time.Unix(0, 0), err
		}
		return time.ParseInLocation(format, ts, loc)
	default:
		return time.Unix(0, 0), errors.New("unsupported type")
	}
}
