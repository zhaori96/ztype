package ztype

import (
	"bytes"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

// Time represents a nullable time value compatible with SQL NULL and JSON null.
// It wraps sql.NullTime with additional JSON parsing capabilities and utility methods.
//
// Example:
//  t := ztype.NewTime(time.Date(2023, time.January, 1, 12, 0, 0, 0, time.UTC))
//  data, _ := json.Marshal(t)
//  // Output: "2023-01-01T12:00:00Z"
type Time struct {
	value       sql.NullTime
	unmarshaled bool
}

var timeFormats = []string{
	time.ANSIC,
	time.UnixDate,
	time.RubyDate,
	time.RFC822,
	time.RFC822Z,
	time.RFC850,
	time.RFC1123,
	time.RFC1123Z,
	time.RFC3339,
	time.RFC3339Nano,
	time.Kitchen,
	time.Stamp,
	time.StampMilli,
	time.StampMicro,
	time.StampNano,
	time.DateTime,
	time.DateOnly,
	time.TimeOnly,
	"2006-01-02 15:04",
	"02/01/2006 15:04:05",
	"02/01/2006 15:04",
	"02/01/2006",
	"15:04",
}

// NewTime creates a non-null Time with an initial value.
//
// Example:
//  t := ztype.NewTime(time.Now())
//  fmt.Println(t.Get().Unix())
func NewTime(value time.Time) Time {
	return Time{
		value: sql.NullTime{
			Time:  value,
			Valid: true,
		},
	}
}

// NewNullTime creates a NULL Time instance.
//
// Example:
//  t := ztype.NewNullTime()
//  fmt.Println(t.IsNull()) // Output: true
func NewNullTime() Time {
	return Time{value: sql.NullTime{Valid: false}}
}

// Get returns the underlying time.Time value.
// Returns zero time if NULL.
//
// Example:
//  value := t.Get()
//  fmt.Println(value.Format(time.RFC822))
func (t *Time) Get() time.Time {
	return t.value.Time
}

// Set updates the value and marks it as valid.
//
// Example:
//  t.Set(time.Date(2025, 5, 15, 0, 0, 0, 0, time.UTC))
func (t *Time) Set(value time.Time) {
	t.value.Time = value
	t.value.Valid = true
}

// SetNull marks the time as NULL.
//
// Example:
//  t.SetNull()
//  fmt.Println(t.IsNull()) // Output: true
func (t *Time) SetNull() {
	t.value.Time = time.Time{}
	t.value.Valid = false
}

// IsNull returns true if the time is NULL.
//
// Example:
//  if t.IsNull() { fmt.Println("Time is NULL") }
func (t *Time) IsNull() bool {
	return !t.value.Valid
}

// IsEmpty returns true if NULL or zero time.
//
// Example:
//  t := ztype.Time{}
//  fmt.Println(t.IsEmpty()) // Output: true
func (t *Time) IsEmpty() bool {
	return !t.value.Valid || t.value.Time.IsZero()
}

// IsZero implements zero value check. Alias for IsEmpty.
//
// Example:
//  t := ztype.NewNullTime()
//  fmt.Println(t.IsZero()) // Output: true
func (t *Time) IsZero() bool {
	return t.IsEmpty()
}

// AddDate adds years, months, and days to the time and returns a new Time.
// Maintains validity state from the original Time.
//
// Example:
//  original := ztype.NewTime(time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC))
//  modified := original.AddDate(1, 2, 3)
//  fmt.Println(modified.Get().Format(time.DateOnly)) // Output: 2024-03-04
func (t Time) AddDate(years int, months int, days int) Time {
	t.value.Time = t.value.Time.AddDate(years, months, days)
	t.value.Valid = true
	return t
}

// AddDateRaw adds years, months, and days to the time and returns the raw time.Time.
// Does not modify the original Time instance.
//
// Example:
//  newTime := t.AddDateRaw(0, 1, 0)
//  fmt.Println(newTime.Month()) // Output: February
func (t *Time) AddDateRaw(years int, months int, days int) time.Time {
	return t.value.Time.AddDate(years, months, days)
}

// Add adds a Duration to the time and returns a new Time.
// Maintains validity state from the original Time.
//
// Example:
//  d := ztype.NewDuration(2 * time.Hour)
//  newTime := t.Add(d)
//  fmt.Println(newTime.Get().Hour())
func (t Time) Add(value Duration) Time {
	t.value.Time = t.value.Time.Add(value.Get())
	t.value.Valid = true
	return t
}

// AddRaw adds a time.Duration to the time and returns the raw time.Time.
// Does not modify the original Time instance.
//
// Example:
//  newTime := t.AddRaw(30 * time.Minute)
//  fmt.Println(newTime.Minute()) // Output: 30
func (t *Time) AddRaw(value time.Duration) time.Time {
	return t.value.Time.Add(value)
}

// Sub calculates duration between two Time values.
// Returns zero Duration if either value is NULL.
//
// Example:
//  duration := t.Sub(otherTime)
//  fmt.Println(duration.Get().Hours())
func (t *Time) Sub(value Time) Duration {
	return NewDuration(t.value.Time.Sub(value.Get()))
}

// SubRaw calculates duration between the Time and a raw time.Time.
//
// Example:
//  diff := t.SubRaw(time.Now())
//  fmt.Println(diff.Seconds())
func (t *Time) SubRaw(value time.Time) time.Duration {
	return t.value.Time.Sub(value)
}

// After reports whether the time is after the given Time.
// Returns false if either value is NULL.
//
// Example:
//  isAfter := t.After(otherTime)
//  fmt.Println(isAfter)
func (t *Time) After(value Time) bool {
	return t.value.Time.After(value.Get())
}

// AfterRaw reports whether the time is after a raw time.Time.
//
// Example:
//  isAfter := t.AfterRaw(time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC))
//  fmt.Println(isAfter)
func (t *Time) AfterRaw(value time.Time) bool {
	return t.value.Time.After(value)
}

// Before reports whether the time is before the given Time.
// Returns false if either value is NULL.
//
// Example:
//  isBefore := t.Before(otherTime)
//  fmt.Println(isBefore)
func (t *Time) Before(value Time) bool {
	return t.value.Time.Before(value.Get())
}

// BeforeRaw reports whether the time is before a raw time.Time.
//
// Example:
//  isBefore := t.BeforeRaw(time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC))
//  fmt.Println(isBefore)
func (t *Time) BeforeRaw(value time.Time) bool {
	return t.value.Time.Before(value)
}

// In returns a copy of the Time with location set to the specified timezone.
//
// Example:
//  loc, _ := time.LoadLocation("America/New_York")
//  nyTime := t.In(loc)
//  fmt.Println(nyTime.Get().Format(time.RFC822))
func (t Time) In(loc *time.Location) Time {
	t.value.Time = t.value.Time.In(loc)
	return t
}

// InRaw returns the raw time.Time in the specified location.
//
// Example:
//  raw := t.InRaw(time.FixedZone("BST", 3600))
//  fmt.Println(raw.Format(time.RFC822Z))
func (t *Time) InRaw(loc *time.Location) time.Time {
	return t.value.Time.In(loc)
}

// Local returns a copy of the Time with location set to the local timezone.
//
// Example:
//  localTime := t.Local()
//  fmt.Println(localTime.Get().Format(time.RFC3339))
func (t Time) Local() Time {
	t.value.Time = t.value.Time.Local()
	return t
}

// LocalRaw returns the raw time.Time in local timezone.
//
// Example:
//  rawLocal := t.LocalRaw()
//  fmt.Println(rawLocal.Format(time.Kitchen))
func (t *Time) LocalRaw() time.Time {
	return t.value.Time.Local()
}

// Location returns the timezone information.
//
// Example:
//  loc := t.Location()
//  fmt.Println(loc.String())
func (t *Time) Location() *time.Location {
	return t.value.Time.Location()
}

// Date returns year, month, and day components.
//
// Example:
//  y, m, d := t.Date()
//  fmt.Printf("%d-%s-%d", y, m, d)
func (t *Time) Date() (year int, month time.Month, day int) {
	return t.value.Time.Date()
}

// Clock returns hour, minute, and second components.
//
// Example:
//  h, m, s := t.Clock()
//  fmt.Printf("%02d:%02d:%02d", h, m, s)
func (t *Time) Clock() (hour, min, sec int) {
	return t.value.Time.Clock()
}

// Nanosecond returns the nanosecond component.
//
// Example:
//  ns := t.Nanosecond()
//  fmt.Println(ns)
func (t *Time) Nanosecond() int {
	return t.value.Time.Nanosecond()
}

// Second returns the second component.
//
// Example:
//  sec := t.Second()
//  fmt.Println(sec)
func (t *Time) Second() int {
	return t.value.Time.Second()
}

// Minute returns the minute component.
//
// Example:
//  min := t.Minute()
//  fmt.Println(min)
func (t *Time) Minute() int {
	return t.value.Time.Minute()
}

// Hour returns the hour component.
//
// Example:
//  hour := t.Hour()
//  fmt.Println(hour)
func (t *Time) Hour() int {
	return t.value.Time.Hour()
}

// Day returns the day component.
//
// Example:
//  day := t.Day()
//  fmt.Println(day)
func (t *Time) Day() int {
	return t.value.Time.Day()
}

// Weekday returns the day of the week.
//
// Example:
//  weekday := t.Weekday()
//  fmt.Println(weekday) // Output: Monday
func (t *Time) Weekday() time.Weekday {
	return t.value.Time.Weekday()
}

// Month returns the month component.
//
// Example:
//  month := t.Month()
//  fmt.Println(month) // Output: January
func (t *Time) Month() time.Month {
	return t.value.Time.Month()
}

// Year returns the year component.
//
// Example:
//  year := t.Year()
//  fmt.Println(year)
func (t *Time) Year() int {
	return t.value.Time.Year()
}

// YearDay returns the day of the year.
//
// Example:
//  yd := t.YearDay()
//  fmt.Println(yd) // Output: 1 for Jan 1
func (t *Time) YearDay() int {
	return t.value.Time.YearDay()
}

// Round returns a new Time rounded to the nearest multiple of the duration.
// Maintains validity state from the original Time.
//
// Example:
//  d := ztype.NewDuration(15 * time.Minute)
//  rounded := t.Round(d)
//  fmt.Println(rounded.Get().Minute()) // Rounds to nearest 15 minutes
func (t Time) Round(value Duration) Time {
	t.value.Time = t.value.Time.Round(value.Get())
	t.value.Valid = true
	return t
}

// RoundRaw rounds the time to the nearest multiple of duration and returns raw time.Time.
//
// Example:
//  raw := t.RoundRaw(time.Hour)
//  fmt.Println(raw.Format(time.TimeOnly))
func (t *Time) RoundRaw(value time.Duration) time.Time {
	return t.value.Time.Round(value)
}

// Truncate returns a new Time truncated to the duration multiple.
// Maintains validity state from the original Time.
//
// Example:
//  d := ztype.NewDuration(24 * time.Hour)
//  truncated := t.Truncate(d)
//  fmt.Println(truncated.Get().Format(time.DateOnly)) // Truncates to midnight
func (t Time) Truncate(value Duration) Time {
	t.value.Time = t.value.Time.Truncate(value.Get())
	t.value.Valid = true
	return t
}

// TruncateRaw truncates the time to duration multiple and returns raw time.Time.
//
// Example:
//  raw := t.TruncateRaw(time.Hour)
//  fmt.Println(raw.Format(time.TimeOnly))
func (t *Time) TruncateRaw(value time.Duration) time.Time {
	return t.value.Time.Truncate(value)
}

// AppendFormat appends formatted time to a byte slice using specified layout.
//
// Example:
//  buf := t.AppendFormat([]byte("Time: "), time.RFC3339)
//  fmt.Println(string(buf))
func (t *Time) AppendFormat(data []byte, layout string) []byte {
	return t.value.Time.AppendFormat(data, layout)
}

// Format returns a string representation using specified layout.
//
// Example:
//  s := t.Format("2006-01-02")
//  fmt.Println(s)
func (t *Time) Format(layout string) string {
	return t.value.Time.Format(layout)
}

// UTC returns a copy of the Time in UTC timezone.
//
// Example:
//  utcTime := t.UTC()
//  fmt.Println(utcTime.Get().Location())
func (t Time) UTC() Time {
	t.value.Time = t.value.Time.UTC()
	return t
}

// UTCRaw returns the raw time.Time in UTC.
//
// Example:
//  utc := t.UTCRaw()
//  fmt.Println(utc.Location())
func (t *Time) UTCRaw() time.Time {
	return t.value.Time.UTC()
}

// Unix returns the Unix timestamp (seconds since epoch).
//
// Example:
//  ts := t.Unix()
//  fmt.Println(ts)
func (t *Time) Unix() int64 {
	return t.value.Time.Unix()
}

// UnixMicro returns the Unix timestamp in microseconds.
//
// Example:
//  ts := t.UnixMicro()
//  fmt.Println(ts)
func (t *Time) UnixMicro() int64 {
	return t.value.Time.UnixMicro()
}

// UnixMilli returns the Unix timestamp in milliseconds.
//
// Example:
//  ts := t.UnixMilli()
//  fmt.Println(ts)
func (t *Time) UnixMilli() int64 {
	return t.value.Time.UnixMilli()
}

// UnixNano returns the Unix timestamp in nanoseconds.
//
// Example:
//  ts := t.UnixNano()
//  fmt.Println(ts)
func (t *Time) UnixNano() int64 {
	return t.value.Time.UnixNano()
}

// GobDecode implements gob.GobDecoder interface.
// Example typically used internally by encoding/gob package.
func (t *Time) GobDecode(data []byte) error {
	return t.value.Time.GobDecode(data)
}

// GobEncode implements gob.GobEncoder interface.
// Example typically used internally by encoding/gob package.
func (t *Time) GobEncode() ([]byte, error) {
	return t.value.Time.GobEncode()
}

// ISOWeek returns the ISO 8601 year and week number.
//
// Example:
//  year, week := t.ISOWeek()
//  fmt.Printf("ISO Week %d of %d", week, year)
func (t *Time) ISOWeek() (year, week int) {
	return t.value.Time.ISOWeek()
}

// Zone returns the timezone name and offset in seconds.
//
// Example:
//  name, offset := t.Zone()
//  fmt.Printf("%s (UTC%+d)", name, offset/3600)
func (t *Time) Zone() (name string, offset int) {
	return t.value.Time.Zone()
}

// Unmarshaled indicates if the value was set through JSON/Text unmarshaling.
//
// Example:
//  if t.Unmarshaled() { fmt.Println("Value from JSON") }
func (t *Time) Unmarshaled() bool {
	return t.unmarshaled
}

// SetUnmarshaled sets the unmarshaled flag status.
// Primarily for internal use.
func (t *Time) SetUnmarshaled(value bool) {
	t.unmarshaled = value
}

// Equal compares both value and null status with another Time.
//
// Example:
//  if t.Equal(otherTime) { fmt.Println("Equal values and null status") }
func (t *Time) Equal(other Time) bool {
	return t.value.Valid == other.value.Valid &&
		t.value.Time.Equal(other.value.Time)
}

// EqualRaw compares the value with a raw time.Time, ignoring null status.
//
// Example:
//  if t.EqualRaw(time.Now()) { fmt.Println("Matches current time") }
func (t *Time) EqualRaw(other time.Time) bool {
	return t.value.Valid && t.value.Time.Equal(other)
}

// MarshalBinary implements encoding.BinaryMarshaler.
// Example typically used internally by encoding packages.
func (t *Time) MarshalBinary() ([]byte, error) {
	return t.value.Time.MarshalBinary()
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
// Example typically used internally by encoding packages.
func (t *Time) UnmarshalBinary(data []byte) error {
	return t.value.Time.UnmarshalBinary(data)
}

// MarshalText implements encoding.TextMarshaler.
// Outputs RFC3339 format for valid times, empty string for NULL.
//
// Example:
//  data, _ := t.MarshalText()
//  fmt.Println(string(data))
func (t *Time) MarshalText() ([]byte, error) {
	if t.value.Valid {
		return []byte(t.value.Time.Format(time.RFC3339)), nil
	}
	return nil, nil
}

// UnmarshalText implements encoding.TextUnmarshaler.
// Supports multiple time formats.
//
// Example:
//  err := t.UnmarshalText([]byte("2023-01-01"))
//  fmt.Println(t.Get().Format(time.DateOnly))
func (t *Time) UnmarshalText(data []byte) error {
	t.unmarshaled = true
	s := string(data)
	if s == "" {
		t.SetNull()
		return nil
	}
	for _, layout := range timeFormats {
		parsed, err := time.Parse(layout, s)
		if err == nil {
			t.value.Time = parsed
			t.value.Valid = true
			return nil
		}
	}
	return fmt.Errorf("invalid time format: %s", s)
}

// MarshalJSON implements json.Marshaler.
// Outputs RFC3339 format for valid times, null for NULL.
//
// Example:
//  data, _ := json.Marshal(t)
//  fmt.Println(string(data))
func (t *Time) MarshalJSON() ([]byte, error) {
	if t.value.Valid {
		return json.Marshal(t.value.Time.Format(time.RFC3339))
	}
	return []byte("null"), nil
}

// UnmarshalJSON implements json.Unmarshaler.
// Supports multiple time formats and null.
//
// Example:
//  err := json.Unmarshal([]byte("\"2023-01-01T00:00:00Z\""), &t)
//  fmt.Println(t.Get().Year())
func (t *Time) UnmarshalJSON(data []byte) error {
	t.unmarshaled = true
	if bytes.Equal(data, []byte("null")) {
		t.SetNull()
		return nil
	}
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	for _, layout := range timeFormats {
		parsed, err := time.Parse(layout, s)
		if err == nil {
			t.value.Time = parsed
			t.value.Valid = true
			return nil
		}
	}
	return fmt.Errorf("invalid time format: %s", s)
}

// Scan implements sql.Scanner for database integration.
//
// Example:
//  err := db.QueryRow("SELECT created_at FROM users").Scan(&t)
func (t *Time) Scan(value any) error {
	return t.value.Scan(value)
}

// Value implements driver.Valuer for database integration.
//
// Example:
//  _, err := db.Exec("INSERT INTO users (created_at) VALUES (?)", t.Value())
func (t *Time) Value() (driver.Value, error) {
	return t.value.Value()
}

// String returns RFC3339Nano format for valid times, "<NULL>" for NULL.
//
// Example:
//  fmt.Println(t.String())
func (t *Time) String() string {
	if !t.value.Valid {
		return "<NULL>"
	}
	return t.value.Time.Format(time.RFC3339Nano)
}

// Duration represents a nullable time.Duration compatible with SQL NULL and JSON null.
//
// Example:
//  d := ztype.NewDuration(5 * time.Minute)
//  data, _ := json.Marshal(d)
//  // Output: "5m0s"
type Duration struct {
	value       time.Duration
	valid       bool
	unmarshaled bool
}

// NewDuration creates a non-null Duration with initial value.
//
// Example:
//  d := ztype.NewDuration(2 * time.Hour)
//  fmt.Println(d.Get().Minutes()) // Output: 120
func NewDuration(value time.Duration) Duration {
	return Duration{
		value: value,
		valid: true,
	}
}

// NewNullDuration creates a NULL Duration instance.
//
// Example:
//  d := ztype.NewNullDuration()
//  fmt.Println(d.IsNull()) // Output: true
func NewNullDuration() Duration {
	return Duration{valid: false}
}

// Get returns the underlying duration value.
// Returns zero duration if NULL.
//
// Example:
//  dur := d.Get()
//  fmt.Println(dur.String())
func (d *Duration) Get() time.Duration {
	return d.value
}

// Set updates the value and marks it as valid.
//
// Example:
//  d.Set(10 * time.Second)
func (d *Duration) Set(value time.Duration) {
	d.value = value
	d.valid = true
}

// SetNull marks the duration as NULL.
//
// Example:
//  d.SetNull()
//  fmt.Println(d.IsNull()) // Output: true
func (d *Duration) SetNull() {
	d.value = 0
	d.valid = false
}

// IsNull returns true if the duration is NULL.
//
// Example:
//  if d.IsNull() { fmt.Println("Duration is NULL") }
func (d *Duration) IsNull() bool {
	return !d.valid
}

// IsZero returns true if NULL or zero duration.
//
// Example:
//  d := ztype.Duration{}
//  fmt.Println(d.IsZero()) // Output: true
func (d *Duration) IsZero() bool {
	return !d.valid || d.value == 0
}

// Unmarshaled indicates if the value was set through JSON/Text unmarshaling.
//
// Example:
//  if d.Unmarshaled() { fmt.Println("Value from JSON") }
func (d *Duration) Unmarshaled() bool {
	return d.unmarshaled
}

// SetUnmarshaled sets the unmarshaled flag status.
// Primarily for internal use.
func (d *Duration) SetUnmarshaled(value bool) {
	d.unmarshaled = value
}

// Equal compares both value and null status with another Duration.
//
// Example:
//  if d.Equal(otherDur) { fmt.Println("Equal values and null status") }
func (d *Duration) Equal(other Duration) bool {
	return d.valid == other.valid && d.value == other.value
}

// EqualRaw compares the value with a raw time.Duration, ignoring null status.
//
// Example:
//  if d.EqualRaw(5 * time.Minute) { fmt.Println("Matches 5 minutes") }
func (d *Duration) EqualRaw(other time.Duration) bool {
	return d.valid && d.value == other
}

// MarshalText implements encoding.TextMarshaler.
// Outputs duration string for valid values, empty string for NULL.
//
// Example:
//  data, _ := d.MarshalText()
//  fmt.Println(string(data))
func (d *Duration) MarshalText() ([]byte, error) {
	if d.valid {
		return []byte(d.value.String()), nil
	}
	return nil, nil
}

// UnmarshalText implements encoding.TextUnmarshaler.
//
// Example:
//  err := d.UnmarshalText([]byte("1h30m"))
//  fmt.Println(d.Get().Minutes()) // Output: 90
func (d *Duration) UnmarshalText(data []byte) error {
	d.unmarshaled = true
	if len(data) == 0 {
		d.SetNull()
		return nil
	}
	dur, err := time.ParseDuration(string(data))
	if err != nil {
		return err
	}
	d.value = dur
	d.valid = true
	return nil
}

// MarshalJSON implements json.Marshaler.
// Outputs duration string for valid values, null for NULL.
//
// Example:
//  data, _ := json.Marshal(d)
//  fmt.Println(string(data)) // Output: "1h30m0s"
func (d *Duration) MarshalJSON() ([]byte, error) {
	if d.valid {
		return json.Marshal(d.value.String())
	}
	return []byte("null"), nil
}

// UnmarshalJSON implements json.Unmarshaler.
//
// Example:
//  err := json.Unmarshal([]byte("\"1h30m\""), &d)
//  fmt.Println(d.Get().Minutes()) // Output: 90
func (d *Duration) UnmarshalJSON(data []byte) error {
	d.unmarshaled = true
	if bytes.Equal(data, []byte("null")) {
		d.SetNull()
		return nil
	}
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	dur, err := time.ParseDuration(s)
	if err != nil {
		return err
	}
	d.value = dur
	d.valid = true
	return nil
}

// Scan implements sql.Scanner for database integration.
// Supports int64 (nanoseconds) and string formats.
//
// Example:
//  err := db.QueryRow("SELECT duration FROM sessions").Scan(&d)
func (d *Duration) Scan(value any) error {
	if value == nil {
		d.value, d.valid = 0, false
		return nil
	}
	switch v := value.(type) {
	case int64:
		d.value = time.Duration(v)
		d.valid = true
	case string:
		dur, err := time.ParseDuration(v)
		if err != nil {
			return err
		}
		d.value = dur
		d.valid = true
	default:
		return fmt.Errorf("unsupported type: %T", value)
	}
	return nil
}

// Value implements driver.Valuer for database integration.
// Returns duration as int64 nanoseconds.
//
// Example:
//  _, err := db.Exec("INSERT INTO sessions (duration) VALUES (?)", d.Value())
func (d Duration) Value() (driver.Value, error) {
	if !d.valid {
		return nil, nil
	}
	return int64(d.value), nil
}

// String returns the duration string for valid values, "<NULL>" for NULL.
//
// Example:
//  fmt.Println(d.String()) // Output: "1h30m0s" or "<NULL>"
func (d *Duration) String() string {
	if !d.valid {
		return "<NULL>"
	}
	return d.value.String()
}
