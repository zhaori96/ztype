package ztype

import (
	"bytes"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"strconv"
)

// Bool represents a nullable boolean type that can distinguish between:
// - Explicit database/SQL NULL values
// - Absent values in JSON unmarshaling
// - Default zero values
//
// It wraps sql.NullBool and adds tracking for unmarshaling presence.
//
// Example Usage:
//
//	// Create valid boolean
//	b := ztype.NewBool(true)
//
//	// Check null state
//	if b.IsNull() { /* handle null case */ }
//
//	// JSON interaction
//	jsonStr := `{"active": true}`
//	err := json.Unmarshal([]byte(jsonStr), &b)
type Bool struct {
	value       sql.NullBool
	unmarshaled bool
}

// NewBool creates a new valid Bool instance.
//
// Example:
//
//	validBool := ztype.NewBool(true)  // Creates non-null true value
//	fmt.Println(validBool.Get())      // Output: true
func NewBool(value bool) Bool {
	return Bool{value: sql.NullBool{Bool: value, Valid: true}}
}

// NewNullBool creates a new null Bool instance.
//
// Example:
//
//	nullBool := ztype.NewNullBool()
//	fmt.Println(nullBool.IsNull())    // Output: true
func NewNullBool() Bool {
	return Bool{value: sql.NullBool{Valid: false}}
}

// NewNullBoolIfZero returns a null Bool if the given value is false.
// Otherwise, it returns a valid Bool with the provided value.
//
// Example:
//
//	b1 := NewNullBoolIfZero(false)   // Null
//	b2 := NewNullBoolIfZero(true)    // Valid with true
func NewNullBoolIfZero(value bool) Bool {
	if !value {
		return NewNullBool()
	}
	return NewBool(value)
}

// Get returns the boolean value. When null, returns false.
// Use IsNull() to check validity before using this value.
//
// Example:
//
//	b := ztype.NewBool(true)
//	if !b.IsNull() {
//	    fmt.Println(b.Get())  // Output: true
//	}
func (b *Bool) Get() bool {
	return b.value.Bool
}

// Set updates the value and marks it as valid.
//
// Example:
//
//	var b ztype.Bool
//	b.Set(true)
//	fmt.Println(b.IsNull())  // Output: false
func (b *Bool) Set(value bool) {
	b.value.Bool = value
	b.value.Valid = true
}

// SetNull marks the value as null and resets the boolean state.
//
// Example:
//
//	b := ztype.NewBool(true)
//	b.SetNull()
//	fmt.Println(b.IsNull())  // Output: true
func (b *Bool) SetNull() {
	b.value.Bool = false
	b.value.Valid = false
}

// IsNull returns true if the value is null.
//
// Example:
//
//	nullBool := ztype.NewNullBool()
//	fmt.Println(nullBool.IsNull())  // Output: true
func (b *Bool) IsNull() bool {
	return !b.value.Valid
}

// IsZero returns true if the value is zero/false.
//
// Example:
//
//	b := ztype.NewBool(false)
//	fmt.Println(b.IsZero())  // Output: true
func (b *Bool) IsZero() bool {
	return !b.value.Bool
}

// Unmarshaled returns true if the value was present in the data source,
// including explicit null values. Returns false if the field was absent.
//
// Example:
//
//	var b ztype.Bool
//	json.Unmarshal([]byte(`{"active": null}`), &b)
//	fmt.Println(b.Unmarshaled())  // Output: true
func (b *Bool) Unmarshaled() bool {
	return b.unmarshaled
}

// SetUnmarshaled manually sets the unmarshaled state. Useful for custom
// serialization/deserialization implementations.
//
// Example:
//
//	b.SetUnmarshaled(true)  // Marks value as coming from external source
func (b *Bool) SetUnmarshaled(value bool) {
	b.unmarshaled = value
}

// Equal performs deep equality check including null state.
//
// Example:
//
//	b1 := ztype.NewBool(true)
//	b2 := ztype.NewBool(true)
//	fmt.Println(b1.Equal(b2))  // Output: true
func (b *Bool) Equal(other Bool) bool {
	return b.value.Bool == other.value.Bool &&
		b.value.Valid == other.value.Valid
}

// EqualRaw compares the boolean value while ignoring null state.
// Returns false if either value is null.
//
// Example:
//
//	b := ztype.NewNullBool()
//	fmt.Println(b.EqualRaw(false))  // Output: false
func (b *Bool) EqualRaw(other bool) bool {
	return b.value.Bool == other
}

// MarshalText implements encoding.TextMarshaler.
// Returns "true"/"false" for valid values, nil for null.
//
// Example:
//
//	b := ztype.NewBool(true)
//	data, _ := b.MarshalText()
//	fmt.Println(string(data))  // Output: true
func (b *Bool) MarshalText() ([]byte, error) {
	if b.value.Valid {
		return []byte(strconv.FormatBool(b.value.Bool)), nil
	}
	return nil, nil
}

// UnmarshalText implements encoding.TextUnmarshaler.
// Sets unmarshaled flag and parses boolean from string.
//
// Example:
//
//	var b ztype.Bool
//	err := b.UnmarshalText([]byte("true"))
//	fmt.Println(b.Get())  // Output: true
func (b *Bool) UnmarshalText(data []byte) error {
	b.unmarshaled = true
	value, err := strconv.ParseBool(string(data))
	if err != nil {
		return err
	}
	b.value.Bool = value
	b.value.Valid = true
	return nil
}

// MarshalJSON implements json.Marshaler.
// Returns JSON boolean for valid values, null for null.
//
// Example:
//
//	b := ztype.NewBool(true)
//	jsonData, _ := json.Marshal(b)
//	fmt.Println(string(jsonData))  // Output: true
func (b *Bool) MarshalJSON() ([]byte, error) {
	if b.value.Valid {
		return json.Marshal(b.value.Bool)
	}
	return []byte("null"), nil
}

// UnmarshalJSON implements json.Unmarshaler.
// Handles both boolean values and explicit nulls.
//
// Example:
//
//	var b ztype.Bool
//	json.Unmarshal([]byte(`null`), &b)
//	fmt.Println(b.IsNull())  // Output: true
func (b *Bool) UnmarshalJSON(data []byte) error {
	b.unmarshaled = true
	if bytes.Equal(data, []byte("null")) {
		b.value.Valid = false
		b.value.Bool = false
		return nil
	}
	b.value.Valid = true
	return json.Unmarshal(data, &b.value.Bool)
}

// Scan implements sql.Scanner for database integration.
//
// Example:
//
//	var b ztype.Bool
//	err := db.QueryRow("SELECT active FROM users WHERE id = 1").Scan(&b)
func (b *Bool) Scan(value any) error {
	return b.value.Scan(value)
}

// Value implements driver.Valuer for database integration.
//
// Example:
//
//	value, _ := b.Value()
//	// Use value in SQL queries
func (b Bool) Value() (driver.Value, error) {
	return b.value.Value()
}

// String returns human-readable representation.
// Returns "<NULL>" for null values, otherwise "true"/"false".
//
// Example:
//
//	b := ztype.NewNullBool()
//	fmt.Println(b.String())  // Output: <NULL>
func (b *Bool) String() string {
	if !b.value.Valid {
		return "<NULL>"
	}
	return strconv.FormatBool(b.value.Bool)
}
