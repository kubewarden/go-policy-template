// This code has been cherry-picked from https://github.com/go-openapi/strfmt

package strfmt

import (
	"time"

	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"
)

type Base64 []byte

type DateTime time.Time

const (
	// RFC3339Millis represents a ISO8601 format to millis instead of to nanos
	RFC3339Millis = "2006-01-02T15:04:05.000Z07:00"
	// RFC3339MillisNoColon represents a ISO8601 format to millis instead of to nanos
	RFC3339MillisNoColon = "2006-01-02T15:04:05.000Z0700"
	// RFC3339Micro represents a ISO8601 format to micro instead of to nano
	RFC3339Micro = "2006-01-02T15:04:05.000000Z07:00"
	// RFC3339MicroNoColon represents a ISO8601 format to micro instead of to nano
	RFC3339MicroNoColon = "2006-01-02T15:04:05.000000Z0700"
	// ISO8601LocalTime represents a ISO8601 format to ISO8601 in local time (no timezone)
	ISO8601LocalTime = "2006-01-02T15:04:05"
	// ISO8601TimeWithReducedPrecision represents a ISO8601 format with reduced precision (dropped secs)
	ISO8601TimeWithReducedPrecision = "2006-01-02T15:04Z"
	// ISO8601TimeWithReducedPrecisionLocaltime represents a ISO8601 format with reduced precision and no timezone (dropped seconds + no timezone)
	ISO8601TimeWithReducedPrecisionLocaltime = "2006-01-02T15:04"
	// ISO8601TimeUniversalSortableDateTimePattern represents a ISO8601 universal sortable date time pattern.
	ISO8601TimeUniversalSortableDateTimePattern = "2006-01-02 15:04:05"

	// json null type
	jsonNull = "null"
)

var (
	// MarshalFormat sets the time resolution format used for marshaling time (set to milliseconds)
	MarshalFormat = RFC3339Millis

	// NormalizeTimeForMarshal provides a normalization function on time befeore marshalling (e.g. time.UTC).
	// By default, the time value is not changed.
	NormalizeTimeForMarshal = func(t time.Time) time.Time { return t }

	// DateTimeFormats is the collection of formats used by ParseDateTime()
	DateTimeFormats = []string{RFC3339Micro, RFC3339MicroNoColon, RFC3339Millis, RFC3339MillisNoColon, time.RFC3339, time.RFC3339Nano, ISO8601LocalTime, ISO8601TimeWithReducedPrecision, ISO8601TimeWithReducedPrecisionLocaltime, ISO8601TimeUniversalSortableDateTimePattern}
)

// MarshalJSON returns the DateTime as JSON
func (t DateTime) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}

	timeForMarshal := NormalizeTimeForMarshal(time.Time(t))

	if timeForMarshal.IsZero() {
		w.RawString(jsonNull)
	} else {
		tstr, err := timeForMarshal.MarshalJSON()
		w.Raw(tstr, err)
	}

	return w.Buffer.BuildBytes(), w.Error
}

// UnmarshalJSON supports json.Unmarshaler interface
func (t *DateTime) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}

	if string(data) == jsonNull {
		return nil
	}

	var tstr string = string(r.String())
	if err := r.Error(); err != nil {
		return err
	}
	tt, err := ParseDateTime(tstr)
	if err != nil {
		return err
	}
	*t = tt
	return nil
}

// ParseDateTime parses a string that represents an ISO8601 time or a unix epoch
func ParseDateTime(data string) (DateTime, error) {
	if data == "" {
		return NewDateTime(), nil
	}
	var lastError error
	for _, layout := range DateTimeFormats {
		dd, err := time.Parse(layout, data)
		if err != nil {
			lastError = err
			continue
		}
		return DateTime(dd), nil
	}
	return DateTime{}, lastError
}

// NewDateTime is a representation of zero value for DateTime type
func NewDateTime() DateTime {
	return DateTime(time.Unix(0, 0).UTC())
}
