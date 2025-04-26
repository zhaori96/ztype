package ztype_test

import (
	"database/sql/driver"
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/zhaori96/ztype"
)

// ============================== Time Tests ==============================

func TestNewTime(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name        string
		input       time.Time
		expectedGet time.Time
		isNull      bool
	}{
		{"valid time", now, now, false},
		{"zero time", time.Time{}, time.Time{}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			zt := ztype.NewTime(tt.input)
			assert.Equal(t, tt.expectedGet, zt.Get())
			assert.Equal(t, tt.isNull, zt.IsNull())
		})
	}
}

func TestNewNullTime(t *testing.T) {
	zt := ztype.NewNullTime()
	assert.True(t, zt.IsNull())
	assert.True(t, zt.IsEmpty())
}

func TestTimeSet(t *testing.T) {
	fixedTime := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)
	newTime := time.Date(2024, 6, 15, 18, 30, 0, 0, time.UTC)
	tests := []struct {
		name         string
		initialValue ztype.Time // Usar ponteiro
		setValue     time.Time
		expectedGet  time.Time
		expectedNull bool
	}{
		{
			name:         "set on null",
			initialValue: ztype.NewNullTime(), // Retorna *ztype.Time
			setValue:     newTime,
			expectedGet:  newTime,
			expectedNull: false,
		},
		{
			name:         "set on valid",
			initialValue: ztype.NewTime(fixedTime), // Retorna *ztype.Time
			setValue:     newTime,
			expectedGet:  newTime,
			expectedNull: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.initialValue.Set(tt.setValue)
			assert.True(t, tt.initialValue.Get().Equal(tt.expectedGet),
				"Expected: %s, Got: %s", tt.expectedGet, tt.initialValue.Get())

			assert.Equal(t, tt.expectedNull, tt.initialValue.IsNull(),
				"Expected null: %v, Got: %v", tt.expectedNull, tt.initialValue.IsNull())
		})
	}
}

func TestTimeSetNull(t *testing.T) {
	zt := ztype.NewTime(time.Now())
	zt.SetNull()
	assert.True(t, zt.IsNull())
	assert.Equal(t, time.Time{}, zt.Get())
}

func TestTimeIsEmpty(t *testing.T) {
	tests := []struct {
		name     string
		input    ztype.Time
		expected bool
	}{
		{"valid non-empty", ztype.NewTime(time.Now()), false},
		{"valid zero", ztype.NewTime(time.Time{}), true},
		{"null", ztype.NewNullTime(), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.input.IsEmpty())
		})
	}
}

func TestTimeAddDate(t *testing.T) {
	original := ztype.NewTime(time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC))
	modified := original.AddDate(1, 2, 3)
	expected := time.Date(2024, 3, 4, 0, 0, 0, 0, time.UTC)
	assert.True(t, modified.Get().Equal(expected))
	assert.False(t, modified.IsNull())
}

func TestTimeAdd(t *testing.T) {
	original := ztype.NewTime(time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC))
	duration := ztype.NewDuration(24 * time.Hour)
	modified := original.Add(duration)
	expected := time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC)
	assert.True(t, modified.Get().Equal(expected))
}

func TestTimeSub(t *testing.T) {
	t1 := ztype.NewTime(time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC))
	t2 := ztype.NewTime(time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC))
	duration := t1.Sub(t2)
	assert.Equal(t, 24*time.Hour, duration.Get())
}

func TestTimeCompare(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name     string
		t1       ztype.Time
		t2       ztype.Time
		expected bool
	}{
		{"after", ztype.NewTime(now.Add(1 * time.Hour)), ztype.NewTime(now), true},
		{"before", ztype.NewTime(now.Add(-1 * time.Hour)), ztype.NewTime(now), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.t1.After(tt.t2))
		})
	}
}

func TestTimeMarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		input    ztype.Time
		expected string
	}{
		{"valid", ztype.NewTime(time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)), `"2023-01-01T12:00:00Z"`},
		{"null", ztype.NewNullTime(), "null"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := tt.input.MarshalJSON()
			assert.NoError(t, err)
			assert.JSONEq(t, tt.expected, string(data))
		})
	}
}

func TestTimeUnmarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		data     string
		expected time.Time
		isNull   bool
	}{
		{"valid", `"2023-01-01T12:00:00Z"`, time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC), false},
		{"null", "null", time.Time{}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var zt ztype.Time
			err := json.Unmarshal([]byte(tt.data), &zt)
			assert.NoError(t, err)
			assert.Equal(t, tt.isNull, zt.IsNull())
			if !tt.isNull {
				assert.True(t, zt.Get().Equal(tt.expected))
			}
		})
	}
}

// ============================== Duration Tests ==============================

func TestNewDuration(t *testing.T) {
	d := ztype.NewDuration(5 * time.Minute)
	assert.Equal(t, 5*time.Minute, d.Get())
	assert.False(t, d.IsNull())
}

func TestDurationSetNull(t *testing.T) {
	d := ztype.NewDuration(10 * time.Second)
	d.SetNull()
	assert.True(t, d.IsNull())
	assert.Equal(t, time.Duration(0), d.Get())
}

func TestDurationMarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		input    ztype.Duration
		expected string
	}{
		{"valid", ztype.NewDuration(2 * time.Hour), `"2h0m0s"`},
		{"null", ztype.NewNullDuration(), "null"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := tt.input.MarshalJSON()
			assert.NoError(t, err)
			assert.JSONEq(t, tt.expected, string(data))
		})
	}
}

func TestDurationUnmarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		data     string
		expected time.Duration
		isNull   bool
	}{
		{"valid", `"1h30m"`, 90 * time.Minute, false},
		{"null", "null", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var d ztype.Duration
			err := json.Unmarshal([]byte(tt.data), &d)
			assert.NoError(t, err)
			assert.Equal(t, tt.isNull, d.IsNull())
			if !tt.isNull {
				assert.Equal(t, tt.expected, d.Get())
			}
		})
	}
}

func TestDurationScan(t *testing.T) {
	tests := []struct {
		name     string
		input    any
		expected time.Duration
		isNull   bool
	}{
		{"int64", int64(3600000000000), time.Hour, false},
		{"string", "1h30m", 90 * time.Minute, false},
		{"null", nil, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var d ztype.Duration
			err := d.Scan(tt.input)
			assert.NoError(t, err)
			assert.Equal(t, tt.isNull, d.IsNull())
			if !tt.isNull {
				assert.Equal(t, tt.expected, d.Get())
			}
		})
	}
}

func TestDurationValue(t *testing.T) {
	tests := []struct {
		name     string
		input    ztype.Duration
		expected driver.Value
	}{
		{"valid", ztype.NewDuration(5 * time.Minute), int64(300000000000)},
		{"null", ztype.NewNullDuration(), nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			val, err := tt.input.Value()
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, val)
		})
	}
}

// ... Adicione mais testes para cobrir todos os métodos restantes
