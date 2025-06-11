package ztype

// Numeric provides a generic nullable numeric type that supports both integer and floating-point types.
// It handles database NULL values, JSON serialization, and arithmetic operations with proper null safety.
//
// Features:
// - Type-safe arithmetic operations
// - Database NULL handling (sql.Scanner/driver.Valuer)
// - JSON serialization/deserialization
// - Null state tracking
// - String formatting
//
// Example usage:
//
//	// Create a valid number
//	num := ztype.NewNumber[int](42)
//
//	// Create a null float
//	nullFloat := ztype.NewNullNumber[float64]()
//
//	// Perform arithmetic
//	sum := num.Add(ztype.NewNumber(10)) // Returns 52
import (
	"bytes"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"math"
	"reflect"
	"strconv"
)

type NumberType interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
		~float32 | ~float64
}

// Numeric represents a nullable numeric value that can be any integer or float type.
// It wraps sql.Null[T] for database compatibility and adds additional functionality.
type Numeric[T NumberType] struct {
	value       sql.Null[T]
	unmarshaled bool
}

// NewNumber creates a new valid Numeric with the specified value.
//
// Example:
//
//	n := NewNumber(42.5)    // Numeric[float64]
//	n := NewNumber(100)     // Numeric[int]
func NewNumber[T NumberType](value T) Numeric[T] {
	return Numeric[T]{value: sql.Null[T]{V: value, Valid: true}}
}

// NewNullNumber creates a new null Numeric of the specified type.
//
// Example:
//
//	n := NewNullNumber[float32]()
func NewNullNumber[T NumberType]() Numeric[T] {
	return Numeric[T]{value: sql.Null[T]{Valid: false}}
}

// NewNullNumberIfZero returns a null Numeric if the given value is zero.
// Otherwise, it returns a valid Numeric with the provided value.
//
// Example:
//
//	n := NewNullNumberIfZero(0)     // Null
//	n2 := NewNullNumberIfZero(42)   // Valid with value 42
func NewNullNumberIfZero[T NumberType](value T) Numeric[T] {
	if value == 0 {
		return NewNullNumber[T]()
	}
	return NewNumber(value)
}

// Get returns the underlying value. Returns zero value if null.
//
// Example:
//
//	n := NewNumber(42)
//	fmt.Println(n.Get()) // Output: 42
func (n *Numeric[T]) Get() T {
	return n.value.V
}

// Set updates the value and marks it as valid.
//
// Example:
//
//	var n Numeric[int]
//	n.Set(100)
//	fmt.Println(n.Get()) // Output: 100
func (n *Numeric[T]) Set(value T) {
	n.value.V = value
	n.value.Valid = true
}

// SetNull marks the value as null and resets the stored value.
//
// Example:
//
//	n.SetNull()
//	fmt.Println(n.IsNull()) // Output: true
func (n *Numeric[T]) SetNull() {
	var zero T
	n.value.V = zero
	n.value.Valid = false
}

// IsNull returns true if the value is null.
//
// Example:
//
//	if num.IsNull() {
//	    fmt.Println("Value is null")
//	}
func (i Numeric[T]) IsNull() bool {
	return !i.value.Valid
}

// Unmarshaled indicates if the value was set through unmarshaling.
// Used for tracking partial updates in data structures.
func (s Numeric[T]) Unmarshaled() bool {
	return s.unmarshaled
}

// SetUnmarshaled controls the unmarshaled flag. Used by parent structures
// to track field state during deserialization.
func (n *Numeric[T]) SetUnmarshaled(value bool) {
	n.unmarshaled = value
}

// Equal compares two Numeric values for equality, including null state.
//
// Example:
//
//	a := NewNumber(10)
//	b := NewNumber(10)
//	fmt.Println(a.Equal(b)) // Output: true
func (n Numeric[T]) Equal(other Numeric[T]) bool {
	return n.value.Valid == other.value.Valid && n.value.V == other.value.V
}

// EqualRaw compares the Numeric value with a raw value.
// Always returns false if the Numeric is null.
//
// Example:
//
//	n := NewNumber(42)
//	fmt.Println(n.EqualRaw(42)) // Output: true
func (n Numeric[T]) EqualRaw(other T) bool {
	return n.value.V == other
}

// Add performs null-safe addition. Returns null if either operand is null.
//
// Example:
//
//	a := NewNumber(10)
//	b := NewNumber(20)
//	c := a.Add(b)
//	fmt.Println(c.Get()) // Output: 30
func (n Numeric[T]) Add(other Numeric[T]) Numeric[T] {
	if !n.value.Valid || !other.value.Valid {
		return NewNullNumber[T]()
	}
	return NewNumber(n.value.V + other.value.V)
}

// AddRaw adds a raw value to the Numeric. Returns zero value if null.
//
// Example:
//
//	n := NewNumber(10)
//	fmt.Println(n.AddRaw(5)) // Output: 15
func (n Numeric[T]) AddRaw(other T) T {
	if n.value.Valid {
		return n.value.V + other
	}
	var zero T
	return zero
}

// Sub performs null-safe subtraction. Returns null if either operand is null.
//
// Example:
//
//	a := NewNumber(30)
//	b := NewNumber(10)
//	c := a.Sub(b)
//	fmt.Println(c.Get()) // Output: 20
func (n Numeric[T]) Sub(other Numeric[T]) Numeric[T] {
	if !n.value.Valid || !other.value.Valid {
		return NewNullNumber[T]()
	}
	return NewNumber(n.value.V - other.value.V)
}

// SubRaw subtracts a raw value from the Numeric. Returns zero value if null.
//
// Example:
//
//	n := NewNumber(20)
//	fmt.Println(n.SubRaw(5)) // Output: 15
func (n Numeric[T]) SubRaw(other T) T {
	if !n.value.Valid {
		var zero T
		return zero
	}
	return n.value.V - other
}

// Mult performs null-safe multiplication. Returns null if either operand is null.
//
// Example:
//
//	a := NewNumber(5)
//	b := NewNumber(4)
//	c := a.Mult(b)
//	fmt.Println(c.Get()) // Output: 20
func (n Numeric[T]) Mult(other Numeric[T]) Numeric[T] {
	if !n.value.Valid || !other.value.Valid {
		return NewNullNumber[T]()
	}
	return NewNumber(n.value.V * other.value.V)
}

// MultRaw multiplies the Numeric by a raw value. Returns zero value if null.
//
// Example:
//
//	n := NewNumber(5)
//	fmt.Println(n.MultRaw(3)) // Output: 15
func (n Numeric[T]) MultRaw(other T) T {
	if !n.value.Valid {
		var zero T
		return zero
	}
	return n.value.V * other
}

// Div performs division. Panics on division by zero or null values.
// Use SafeDiv for error handling version.
//
// Example:
//
//	a := NewNumber(20)
//	b := NewNumber(5)
//	c := a.Div(b)
//	fmt.Println(c.Get()) // Output: 4
func (n Numeric[T]) Div(other Numeric[T]) Numeric[T] {
	value, err := n.SafeDiv(other)
	if err != nil {
		panic(err)
	}
	return value
}

// SafeDiv performs null-safe division with error handling.
// Returns error for division by zero or null values.
//
// Example:
//
//	a := NewNumber(20)
//	b := NewNumber(0)
//	_, err := a.SafeDiv(b)
//	fmt.Println(err) // Output: cannot divide by zero
func (n Numeric[T]) SafeDiv(other Numeric[T]) (Numeric[T], error) {
	if !other.value.Valid || other.value.V == 0 {
		return NewNullNumber[T](), fmt.Errorf("cannot divide by zero")
	}
	return NewNumber(n.value.V / other.value.V), nil
}

// DivRaw divides by a raw value. Panics on division by zero.
//
// Example:
//
//	n := NewNumber(20)
//	fmt.Println(n.DivRaw(5)) // Output: 4
func (n Numeric[T]) DivRaw(other T) T {
	value, err := n.SafeDivRaw(other)
	if err != nil {
		panic(err)
	}
	return value
}

// SafeDivRaw divides by a raw value with error handling.
//
// Example:
//
//	n := NewNumber(20)
//	result, err := n.SafeDivRaw(0)
//	fmt.Println(err) // Output: cannot divide by zero
func (n Numeric[T]) SafeDivRaw(other T) (T, error) {
	if other == 0 {
		return 0, fmt.Errorf("cannot divide by zero")
	}
	return n.value.V / other, nil
}

// Compare compares two Numeric values. Returns:
// -1 if n < other
//
//	0 if n == other
//	1 if n > other
//
// Error if either value is null.
//
// Example:
//
//	a := NewNumber(10)
//	b := NewNumber(20)
//	result, _ := a.Compare(b)
//	fmt.Println(result) // Output: -1
func (n Numeric[T]) Compare(other Numeric[T]) (int, error) {
	if !n.value.Valid || !other.value.Valid {
		return 0, fmt.Errorf("cannot compare null values")
	}
	if n.value.V < other.value.V {
		return -1, nil
	} else if n.value.V > other.value.V {
		return 1, nil
	}
	return 0, nil
}

// CompareRaw compares with a raw value. Returns error if null.
//
// Example:
//
//	n := NewNumber(42)
//	result, _ := n.CompareRaw(30)
//	fmt.Println(result) // Output: 1
func (n Numeric[T]) CompareRaw(other T) (int, error) {
	if !n.value.Valid {
		return 0, fmt.Errorf("cannot compare null values")
	}
	if n.value.V < other {
		return -1, nil
	} else if n.value.V > other {
		return 1, nil
	}
	return 0, nil
}

// Greater returns true if n > other. Returns false if either is null.
//
// Example:
//
//	a := NewNumber(20)
//	b := NewNumber(10)
//	fmt.Println(a.Greater(b)) // Output: true
func (n Numeric[T]) Greater(other Numeric[T]) bool {
	if !n.value.Valid || !other.value.Valid {
		return false
	}
	return n.value.V > other.value.V
}

// GreaterRaw returns true if n > raw value. Returns false if null.
//
// Example:
//
//	n := NewNumber(15)
//	fmt.Println(n.GreaterRaw(10)) // Output: true
func (n Numeric[T]) GreaterRaw(other T) bool {
	if !n.value.Valid {
		return false
	}
	return n.value.V > other
}

// GreaterOrEqual returns true if n >= other. Returns false if either is null.
//
// Example:
//
//	a := NewNumber(10)
//	b := NewNumber(10)
//	fmt.Println(a.GreaterOrEqual(b)) // Output: true
func (n Numeric[T]) GreaterOrEqual(other Numeric[T]) bool {
	if !n.value.Valid || !other.value.Valid {
		return false
	}
	return n.value.V >= other.value.V
}

// GreaterOrEqualRaw returns true if n >= raw value. Returns false if null.
//
// Example:
//
//	n := NewNumber(10)
//	fmt.Println(n.GreaterOrEqualRaw(10)) // Output: true
func (n Numeric[T]) GreaterOrEqualRaw(other T) bool {
	if !n.value.Valid {
		return false
	}
	return n.value.V >= other
}

// Less returns true if n < other. Returns false if either is null.
//
// Example:
//
//	a := NewNumber(5)
//	b := NewNumber(10)
//	fmt.Println(a.Less(b)) // Output: true
func (n Numeric[T]) Less(other Numeric[T]) bool {
	if !n.value.Valid || !other.value.Valid {
		return false
	}
	return n.value.V < other.value.V
}

// LessRaw returns true if n < raw value. Returns false if null.
//
// Example:
//
//	n := NewNumber(5)
//	fmt.Println(n.LessRaw(10)) // Output: true
func (n Numeric[T]) LessRaw(other T) bool {
	if !n.value.Valid {
		return false
	}
	return n.value.V < other
}

// LessOrEqual returns true if n <= other. Returns false if either is null.
//
// Example:
//
//	a := NewNumber(10)
//	b := NewNumber(10)
//	fmt.Println(a.LessOrEqual(b)) // Output: true
func (n Numeric[T]) LessOrEqual(other Numeric[T]) bool {
	if !n.value.Valid || !other.value.Valid {
		return false
	}
	return n.value.V <= other.value.V
}

// LessOrEqualRaw returns true if n <= raw value. Returns false if null.
//
// Example:
//
//	n := NewNumber(10)
//	fmt.Println(n.LessOrEqualRaw(10)) // Output: true
func (n Numeric[T]) LessOrEqualRaw(other T) bool {
	if !n.value.Valid {
		return false
	}
	return n.value.V <= other
}

// Min returns the smaller of two Numeric values. Treats null as negative infinity.
//
// Example:
//
//	a := NewNumber(5)
//	b := NewNumber(10)
//	fmt.Println(a.Min(b).Get()) // Output: 5
func (n Numeric[T]) Min(other Numeric[T]) Numeric[T] {
	if !n.value.Valid && !other.value.Valid {
		return NewNullNumber[T]()
	}
	if !n.value.Valid {
		return other
	}
	if !other.value.Valid {
		return n
	}
	if n.value.V <= other.value.V {
		return n
	}
	return other
}

// MinRaw returns the smaller of the Numeric value and a raw value.
//
// Example:
//
//	n := NewNumber(5)
//	fmt.Println(n.MinRaw(3)) // Output: 3
func (n Numeric[T]) MinRaw(other T) T {
	if !n.value.Valid {
		return other
	}
	if n.value.V <= other {
		return n.value.V
	}
	return other
}

// Max returns the larger of two Numeric values. Treats null as positive infinity.
//
// Example:
//
//	a := NewNumber(5)
//	b := NewNumber(10)
//	fmt.Println(a.Max(b).Get()) // Output: 10
func (n Numeric[T]) Max(other Numeric[T]) Numeric[T] {
	if !n.value.Valid && !other.value.Valid {
		return NewNullNumber[T]()
	}
	if !n.value.Valid {
		return other
	}
	if !other.value.Valid {
		return n
	}
	if n.value.V >= other.value.V {
		return n
	}
	return other
}

// MaxRaw returns the larger of the Numeric value and a raw value.
//
// Example:
//
//	n := NewNumber(5)
//	fmt.Println(n.MaxRaw(10)) // Output: 10
func (n Numeric[T]) MaxRaw(other T) T {
	if !n.value.Valid {
		return other
	}
	if n.value.V >= other {
		return n.value.V
	}
	return other
}

// MarshalText implements encoding.TextMarshaler.
//
// Example:
//
//	n := NewNumber(123.456)
//	data, _ := n.MarshalText()
//	fmt.Println(string(data)) // Output: 123.456000
func (n *Numeric[T]) MarshalText() ([]byte, error) {
	if n.value.Valid {
		return []byte(n.String()), nil
	}
	return nil, nil
}

// UnmarshalText implements encoding.TextUnmarshaler.
//
// Example:
//
//	var n Numeric[float64]
//	n.UnmarshalText([]byte("123.45"))
//	fmt.Println(n.Get()) // Output: 123.45
func (n *Numeric[T]) UnmarshalText(data []byte) error {
	n.unmarshaled = true
	if len(data) == 0 {
		n.value.Valid = false
		return nil
	}

	var value T
	var kind reflect.Kind = reflect.TypeOf(value).Kind()
	switch kind {
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		parsed, err := parseUint[T](data, kind)
		if err != nil {
			return err
		}
		value = parsed
	case reflect.Float32, reflect.Float64:
		parsed, err := parseFloat[T](data, kind)
		if err != nil {
			return err
		}
		value = parsed
	default:
		parsed, err := parseInt[T](data, kind)
		if err != nil {
			return err
		}
		value = T(parsed)
	}

	n.value.V = value
	n.value.Valid = true
	return nil
}

// MarshalJSON implements json.Marshaler.
//
// Example:
//
//	n := NewNumber(3.14)
//	j, _ := json.Marshal(n)
//	fmt.Println(string(j)) // Output: 3.14
func (n *Numeric[T]) MarshalJSON() ([]byte, error) {
	if n.value.Valid {
		return json.Marshal(n.value.V)
	}
	return []byte("null"), nil
}

// UnmarshalJSON implements json.Unmarshaler.
//
// Example:
//
//	var n Numeric[int]
//	json.Unmarshal([]byte("100"), &n)
//	fmt.Println(n.Get()) // Output: 100
func (n *Numeric[T]) UnmarshalJSON(data []byte) error {
	n.unmarshaled = true
	if bytes.Equal(data, []byte("null")) {
		var zero T
		n.value.Valid = false
		n.value.V = zero
		return nil
	}

	var value T
	if err := json.Unmarshal(data, &value); err != nil {
		n.value.Valid = false
		return err
	}

	n.value.Valid = true
	n.value.V = value
	return nil
}

// Scan implements sql.Scanner for database operations.
//
// Example:
//
//	var n Numeric[float64]
//	db.QueryRow("SELECT price FROM products").Scan(&n)
func (n *Numeric[T]) Scan(value any) error {
	return n.value.Scan(value)
}

// Value implements driver.Valuer for database operations.
//
// Example:
//
//	n := NewNumber(42)
//	val, _ := n.Value()
//	fmt.Printf("%T", val) // Output: int
func (n Numeric[T]) Value() (driver.Value, error) {
	return n.value.Value()
}

// String returns a human-readable representation.
//
// Example:
//
//	n := NewNumber(123.456)
//	fmt.Println(n.String()) // Output: 123.456000
func (n *Numeric[T]) String() string {
	if !n.value.Valid {
		return "<NULL>"
	}

	switch value := any(n.value.V).(type) {
	case float32, float64:
		return fmt.Sprintf("%f", value)
	default:
		return fmt.Sprintf("%v", value)
	}
}

// parseFloat converts byte data to float types with overflow checking.
func parseFloat[T NumberType](
	data []byte,
	kind reflect.Kind,
) (T, error) {
	var zero T
	parsed, err := strconv.ParseFloat(string(data), 64)
	if err != nil {
		return zero, err
	}

	if kind == reflect.Float32 && (parsed > math.MaxFloat32 || parsed < -math.MaxFloat32) {
		return zero, fmt.Errorf("value %f overflows float32", parsed)
	}
	return T(parsed), nil
}

// parseUint converts byte data to unsigned integer types with overflow checking.
func parseUint[T NumberType](
	data []byte,
	kind reflect.Kind,
) (T, error) {
	var zero T
	parsed, err := strconv.ParseUint(string(data), 10, 64)
	if err != nil {
		return zero, err
	}

	switch kind {
	case reflect.Uint:
		if parsed > math.MaxUint {
			return zero, fmt.Errorf("value %d overflows uint", parsed)
		}
	case reflect.Uint8:
		if parsed > math.MaxUint8 {
			return zero, fmt.Errorf("value %d overflows uint8", parsed)
		}
	case reflect.Uint16:
		if parsed > math.MaxUint16 {
			return zero, fmt.Errorf("value %d overflows uint16", parsed)
		}
	case reflect.Uint32:
		if parsed > math.MaxUint32 {
			return zero, fmt.Errorf("value %d overflows uint32", parsed)
		}
	}

	return T(parsed), nil
}

// parseInt converts byte data to signed integer types with overflow checking.
func parseInt[T NumberType](
	data []byte,
	kind reflect.Kind,
) (T, error) {
	var zero T
	parsed, err := strconv.ParseInt(string(data), 10, 64)
	if err != nil {
		return zero, err
	}

	switch kind {
	case reflect.Int:
		if parsed > math.MaxInt || parsed < math.MinInt {
			return zero, fmt.Errorf("value %d overflows int", parsed)
		}
	case reflect.Int8:
		if parsed > math.MaxInt8 || parsed < math.MinInt8 {
			return zero, fmt.Errorf("value %d overflows int8", parsed)
		}
	case reflect.Int16:
		if parsed > math.MaxInt16 || parsed < math.MinInt16 {
			return zero, fmt.Errorf("value %d overflows int16", parsed)
		}
	case reflect.Int32:
		if parsed > math.MaxInt32 || parsed < math.MinInt32 {
			return zero, fmt.Errorf("value %d overflows int32", parsed)
		}
	}

	return T(parsed), nil
}
