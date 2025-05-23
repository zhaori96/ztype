package ztype

import (
	"bytes"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
)

// String represents a nullable string compatible with SQL NULL and JSON null.
//
// Example declarations:
//
//	var s1 ztype.String = ztype.NewString("text")
//	s2 := ztype.NewNullString()
type String struct {
	value       sql.NullString
	unmarshaled bool
}

// NewString creates a non-null String with initial value.
//
// Example:
//
//	s := ztype.NewString("initial")
//	s.Get() // returns "initial"
//	s.IsNull() // false
func NewString(value string) String {
	return String{value: sql.NullString{String: value, Valid: true}}
}

// NewNullString creates a NULL String instance.
//
// Example:
//
//	s := ztype.NewNullString()
//	s.IsNull() // true
func NewNullString() String {
	return String{value: sql.NullString{Valid: false}}
}

// Get returns the underlying string value (empty if NULL).
//
// Example:
//
//	s := ztype.NewString("value")
//	s.Get() // "value"
func (s *String) Get() string {
	return s.value.String
}

// Set updates the string value and marks it as valid.
//
// Example:
//
//	var s ztype.String
//	s.Set("new-value")
//	s.Get() // "new-value"
//	s.IsNull() // false
func (s *String) Set(value string) {
	s.value.String = value
	s.value.Valid = true
}

// SetNull marks the string as NULL.
//
// Example:
//
//	s := ztype.NewString("text")
//	s.SetNull()
//	s.IsNull() // true
func (s *String) SetNull() {
	s.value.String = ""
	s.value.Valid = false
}

// IsNull returns true if the string is NULL.
//
// Example:
//
//	s := ztype.NewNullString()
//	s.IsNull() // true
func (s *String) IsNull() bool {
	return !s.value.Valid
}

// IsEmpty returns true if NULL or empty string.
//
// Example:
//
//	s1 := ztype.NewNullString()
//	s2 := ztype.NewString("")
//	s1.IsEmpty() // true
//	s2.IsEmpty() // true
func (s *String) IsEmpty() bool {
	return !s.value.Valid || s.value.String == ""
}

// IsZero implements common interface for zero checks (alias for IsEmpty).
//
// Example:
//
//	s := ztype.NewString("")
//	s.IsZero() // true
func (s *String) IsZero() bool {
	return s.IsEmpty()
}

// Unmarshaled indicates if value was set via JSON/text unmarshaling.
//
// Example:
//
//	var s ztype.String
//	json.Unmarshal([]byte(`"data"`), &s)
//	s.Unmarshaled() // true
func (s *String) Unmarshaled() bool {
	return s.unmarshaled
}

// SetUnmarshaled manually controls the unmarshaled flag.
//
// Example:
//
//	s := ztype.NewString("test")
//	s.SetUnmarshaled(true)
func (s *String) SetUnmarshaled(value bool) {
	s.unmarshaled = value
}

// Equal compares both value and null state of two Strings.
//
// Example:
//
//	s1 := ztype.NewString("a")
//	s2 := ztype.NewString("a")
//	s3 := ztype.NewNullString()
//	s1.Equal(s2) // true
//	s1.Equal(s3) // false
func (s *String) Equal(other String) bool {
	return s.value.String == other.value.String &&
		s.value.Valid == other.value.Valid
}

// EqualRaw compares value ignoring null state.
//
// Example:
//
//	s := ztype.NewString("test")
//	s.EqualRaw("test") // true
//	s.EqualRaw("other") // false
func (s *String) EqualRaw(other string) bool {
	return s.value.String == other
}

// MarshalText implements encoding.TextMarshaler.
//
// Example:
//
//	s := ztype.NewString("text")
//	data, _ := s.MarshalText()
//	string(data) // "text"
func (s *String) MarshalText() ([]byte, error) {
	if s.value.Valid {
		return []byte(s.value.String), nil
	}
	return nil, nil
}

// UnmarshalText implements encoding.TextUnmarshaler.
//
// Example:
//
//	var s ztype.String
//	s.UnmarshalText([]byte("data"))
//	s.Get() // "data"
//	s.Unmarshaled() // true
func (s *String) UnmarshalText(data []byte) error {
	s.unmarshaled = true
	s.value.String = string(data)
	s.value.Valid = true
	return nil
}

// MarshalJSON implements json.Marshaler.
//
// Example:
//
//	s := ztype.NewNullString()
//	data, _ := json.Marshal(s)
//	string(data) // "null"
func (s *String) MarshalJSON() ([]byte, error) {
	if s.value.Valid {
		return json.Marshal(s.value.String)
	}
	return []byte("null"), nil
}

// UnmarshalJSON implements json.Unmarshaler.
//
// Example:
//
//	var s ztype.String
//	json.Unmarshal([]byte(`"json-value"`), &s)
//	s.Get() // "json-value"
func (s *String) UnmarshalJSON(data []byte) error {
	s.unmarshaled = true
	if bytes.Equal(data, []byte("null")) {
		s.value.Valid = false
		s.value.String = ""
		return nil
	}
	s.value.Valid = true
	return json.Unmarshal(data, &s.value.String)
}

// Scan implements sql.Scanner for database integration.
//
// Example:
//
//	var s ztype.String
//	s.Scan("scanned-value")
//	s.Get() // "scanned-value"
func (s *String) Scan(value any) error {
	return s.value.Scan(value)
}

// Value implements driver.Valuer for database integration.
//
// Example:
//
//	s := ztype.NewString("db-value")
//	val, _ := s.Value()
//	val.(string) // "db-value"
func (s String) Value() (driver.Value, error) {
	return s.value.Value()
}

// String implements fmt.Stringer for human-readable output.
//
// Example:
//
//	s := ztype.NewNullString()
//	fmt.Println(s) // "<NULL>"
func (s *String) String() string {
	if !s.value.Valid {
		return "<NULL>"
	}
	return s.value.String
}
