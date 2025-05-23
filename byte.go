package ztype

import (
	"bytes"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"strconv"
)

// Byte represents a nullable byte type that can distinguish between:
// - Explicit database/SQL NULL values
// - Absent values in JSON unmarshaling
// - Default zero values
//
// It wraps sql.NullByte and adds tracking for unmarshaling presence.
//
// Example Usage:
//
//	// Create valid byte
//	b := ztype.NewByte(10)
//
//	// Check null state
//	if b.IsNull() { /* handle null case */ }
//
//	// JSON interaction
//	jsonStr := `{"value": 5}`
//	err := json.Unmarshal([]byte(jsonStr), &b)
type Byte struct {
	value       sql.NullByte
	unmarshaled bool
}

// NewByte creates a new valid Byte instance.
//
// Example:
//
//	validByte := ztype.NewByte(10)  // Creates non-null byte 10
//	fmt.Println(validByte.Get())    // Output: 10
func NewByte(value byte) Byte {
	return Byte{value: sql.NullByte{Byte: value, Valid: true}}
}

// NewNullByte creates a new null Byte instance.
//
// Example:
//
//	nullByte := ztype.NewNullByte()
//	fmt.Println(nullByte.IsNull())  // Output: true
func NewNullByte() Byte {
	return Byte{value: sql.NullByte{Valid: false}}
}

// Get returns the byte value. When null, returns 0.
// Use IsNull() to check validity before using this value.
//
// Example:
//
//	b := ztype.NewByte(5)
//	if !b.IsNull() {
//	    fmt.Println(b.Get())  // Output: 5
//	}
func (b *Byte) Get() byte {
	return b.value.Byte
}

// Set updates the value and marks it as valid.
//
// Example:
//
//	var b ztype.Byte
//	b.Set(10)
//	fmt.Println(b.IsNull())  // Output: false
func (b *Byte) Set(value byte) {
	b.value.Byte = value
	b.value.Valid = true
}

// SetNull marks the value as null and resets the byte state.
//
// Example:
//
//	b := ztype.NewByte(5)
//	b.SetNull()
//	fmt.Println(b.IsNull())  // Output: true
func (b *Byte) SetNull() {
	b.value.Byte = 0
	b.value.Valid = false
}

// IsNull returns true if the value is null.
//
// Example:
//
//	nullByte := ztype.NewNullByte()
//	fmt.Println(nullByte.IsNull())  // Output: true
func (b *Byte) IsNull() bool {
	return !b.value.Valid
}

// IsZero returns true if the value is zero.
//
// Example:
//
//	b := ztype.NewByte(0)
//	fmt.Println(b.IsZero())  // Output: true
func (b *Byte) IsZero() bool {
	return !b.value.Valid
}

// Unmarshaled returns true if the value was present in the data source,
// including explicit null values. Returns false if the field was absent.
//
// Example:
//
//	var b ztype.Byte
//	json.Unmarshal([]byte(`{"value": null}`), &b)
//	fmt.Println(b.Unmarshaled())  // Output: true
func (b *Byte) Unmarshaled() bool {
	return b.unmarshaled
}

// SetUnmarshaled manually sets the unmarshaled state. Useful for custom
// serialization/deserialization implementations.
//
// Example:
//
//	b.SetUnmarshaled(true)  // Marks value as coming from external source
func (b *Byte) SetUnmarshaled(value bool) {
	b.unmarshaled = value
}

// Equal performs deep equality check including null state.
//
// Example:
//
//	b1 := ztype.NewByte(5)
//	b2 := ztype.NewByte(5)
//	fmt.Println(b1.Equal(b2))  // Output: true
func (b *Byte) Equal(other Byte) bool {
	return b.value.Byte == other.value.Byte &&
		b.value.Valid == other.value.Valid
}

// EqualRaw compares the byte value while ignoring null state.
// Returns false if either value is null.
//
// Example:
//
//	b := ztype.NewNullByte()
//	fmt.Println(b.EqualRaw(0))  // Output: false
func (b *Byte) EqualRaw(other byte) bool {
	return b.value.Byte == other
}

// MarshalText implements encoding.TextMarshaler.
// Returns string representation for valid values, nil for null.
//
// Example:
//
//	b := ztype.NewByte(10)
//	data, _ := b.MarshalText()
//	fmt.Println(string(data))  // Output: 10
func (b *Byte) MarshalText() ([]byte, error) {
	if b.value.Valid {
		return []byte(strconv.FormatUint(uint64(b.value.Byte), 10)), nil
	}
	return nil, nil
}

// UnmarshalText implements encoding.TextUnmarshaler.
// Sets unmarshaled flag and parses byte from string.
//
// Example:
//
//	var b ztype.Byte
//	err := b.UnmarshalText([]byte("255"))
//	fmt.Println(b.Get())  // Output: 255
func (b *Byte) UnmarshalText(data []byte) error {
	b.unmarshaled = true
	value, err := strconv.ParseUint(string(data), 10, 8)
	if err != nil {
		return err
	}
	b.value.Byte = byte(value)
	b.value.Valid = true
	return nil
}

// MarshalJSON implements json.Marshaler.
// Returns JSON number for valid values, null for null.
//
// Example:
//
//	b := ztype.NewByte(10)
//	jsonData, _ := json.Marshal(b)
//	fmt.Println(string(jsonData))  // Output: 10
func (b *Byte) MarshalJSON() ([]byte, error) {
	if b.value.Valid {
		return json.Marshal(b.value.Byte)
	}
	return []byte("null"), nil
}

// UnmarshalJSON implements json.Unmarshaler.
// Handles both numeric values and explicit nulls.
//
// Example:
//
//	var b ztype.Byte
//	json.Unmarshal([]byte(`null`), &b)
//	fmt.Println(b.IsNull())  // Output: true
func (b *Byte) UnmarshalJSON(data []byte) error {
	b.unmarshaled = true
	if bytes.Equal(data, []byte("null")) {
		b.value.Valid = false
		b.value.Byte = 0
		return nil
	}
	if err := json.Unmarshal(data, &b.value.Byte); err != nil {
		b.value.Valid = false
		return err
	}
	b.value.Valid = true
	return nil
}

// Scan implements sql.Scanner for database integration.
//
// Example:
//
//	var b ztype.Byte
//	err := db.QueryRow("SELECT value FROM table WHERE id = 1").Scan(&b)
func (b *Byte) Scan(value any) error {
	return b.value.Scan(value)
}

// Value implements driver.Valuer for database integration.
//
// Example:
//
//	value, _ := b.Value()
//	// Use value in SQL queries
func (b Byte) Value() (driver.Value, error) {
	return b.value.Value()
}

// String returns human-readable representation.
// Returns "<NULL>" for null values, decimal string otherwise.
//
// Example:
//
//	b := ztype.NewNullByte()
//	fmt.Println(b.String())  // Output: <NULL>
func (b *Byte) String() string {
	if !b.value.Valid {
		return "<NULL>"
	}
	return strconv.FormatUint(uint64(b.value.Byte), 10)
}
