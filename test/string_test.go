package ztype_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/zhaori96/ztype"
)

func TestNewString(t *testing.T) {
	tests := []struct {
		name            string
		input           string
		expectedGet     string
		expectedIsNull  bool
		expectedIsEmpty bool
	}{
		{
			name:            "non-empty value",
			input:           "text",
			expectedGet:     "text",
			expectedIsNull:  false,
			expectedIsEmpty: false,
		},
		{
			name:            "empty value",
			input:           "",
			expectedGet:     "",
			expectedIsNull:  false,
			expectedIsEmpty: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := ztype.NewString(tt.input)
			assert.Equal(t, tt.expectedGet, s.Get())
			assert.Equal(t, tt.expectedIsNull, s.IsNull())
			assert.Equal(t, tt.expectedIsEmpty, s.IsEmpty())
		})
	}
}

func TestNewNullString(t *testing.T) {
	s := ztype.NewNullString()
	assert.True(t, s.IsNull())
	assert.True(t, s.IsEmpty())
	assert.Equal(t, "", s.Get())
}

func TestSet(t *testing.T) {
	tests := []struct {
		name         string
		initialValue ztype.String
		setValue     string
		expectedGet  string
		expectedNull bool
	}{
		{
			name:         "set non-empty on null",
			initialValue: ztype.NewNullString(),
			setValue:     "new",
			expectedGet:  "new",
			expectedNull: false,
		},
		{
			name:         "set empty on non-null",
			initialValue: ztype.NewString("old"),
			setValue:     "",
			expectedGet:  "",
			expectedNull: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.initialValue.Set(tt.setValue)
			assert.Equal(t, tt.expectedGet, tt.initialValue.Get())
			assert.Equal(t, tt.expectedNull, tt.initialValue.IsNull())
		})
	}
}

func TestSetNull(t *testing.T) {
	s := ztype.NewString("text")
	s.SetNull()
	assert.True(t, s.IsNull())
	assert.Equal(t, "", s.Get())
}

func TestIsNull(t *testing.T) {
	tests := []struct {
		name     string
		input    ztype.String
		expected bool
	}{
		{"non-null", ztype.NewString("text"), false},
		{"null", ztype.NewNullString(), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.input.IsNull())
		})
	}
}

func TestIsEmpty(t *testing.T) {
	tests := []struct {
		name     string
		input    ztype.String
		expected bool
	}{
		{"non-null non-empty", ztype.NewString("text"), false},
		{"non-null empty", ztype.NewString(""), true},
		{"null", ztype.NewNullString(), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.input.IsEmpty())
		})
	}
}

func TestIsZero(t *testing.T) {
	tests := []struct {
		name     string
		input    ztype.String
		expected bool
	}{
		{"non-zero", ztype.NewString("text"), false},
		{"zero (empty)", ztype.NewString(""), true},
		{"zero (null)", ztype.NewNullString(), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.input.IsZero())
		})
	}
}

func TestUnmarshaled(t *testing.T) {
	var s ztype.String
	assert.False(t, s.Unmarshaled())

	_ = json.Unmarshal([]byte(`"data"`), &s)
	assert.True(t, s.Unmarshaled())
}

func TestSetUnmarshaled(t *testing.T) {
	s := ztype.NewString("test")
	s.SetUnmarshaled(true)
	assert.True(t, s.Unmarshaled())
	s.SetUnmarshaled(false)
	assert.False(t, s.Unmarshaled())
}

func TestEqual(t *testing.T) {
	tests := []struct {
		name     string
		a        ztype.String
		b        ztype.String
		expected bool
	}{
		{"both non-null equal", ztype.NewString("a"), ztype.NewString("a"), true},
		{"non-null different", ztype.NewString("a"), ztype.NewString("b"), false},
		{"a null, b non-null", ztype.NewNullString(), ztype.NewString("a"), false},
		{"both null", ztype.NewNullString(), ztype.NewNullString(), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.a.Equal(tt.b))
		})
	}
}

func TestEqualRaw(t *testing.T) {
	tests := []struct {
		name     string
		input    ztype.String
		compare  string
		expected bool
	}{
		{"non-null equal", ztype.NewString("a"), "a", true},
		{"non-null different", ztype.NewString("a"), "b", false},
		{"null compare empty", ztype.NewNullString(), "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.input.EqualRaw(tt.compare))
		})
	}
}

func TestMarshalText(t *testing.T) {
	tests := []struct {
		name     string
		input    ztype.String
		expected []byte
	}{
		{"non-null", ztype.NewString("text"), []byte("text")},
		{"null", ztype.NewNullString(), nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := tt.input.MarshalText()
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, data)
		})
	}
}

func TestUnmarshalText(t *testing.T) {
	tests := []struct {
		name     string
		data     []byte
		expected ztype.String
	}{
		{"non-empty", []byte("data"), ztype.NewString("data")},
		{"empty", []byte(""), ztype.NewString("")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var s ztype.String
			err := s.UnmarshalText(tt.data)
			assert.NoError(t, err)
			assert.True(t, s.Unmarshaled())
			assert.Equal(t, tt.expected.Get(), s.Get())
			assert.Equal(t, tt.expected.IsNull(), s.IsNull())
		})
	}
}

func TestMarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		input    ztype.String
		expected []byte
	}{
		{"non-null", ztype.NewString("text"), []byte(`"text"`)},
		{"null", ztype.NewNullString(), []byte("null")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := tt.input.MarshalJSON()
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, data)
		})
	}
}

func TestUnmarshalJSON(t *testing.T) {
	tests := []struct {
		name        string
		data        []byte
		expected    ztype.String
		expectError bool
	}{
		{"valid string", []byte(`"text"`), ztype.NewString("text"), false},
		{"null", []byte("null"), ztype.NewNullString(), false},
		{"invalid type", []byte("123"), ztype.String{}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var s ztype.String
			err := json.Unmarshal(tt.data, &s)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.True(t, s.Unmarshaled())
				assert.Equal(t, tt.expected.Get(), s.Get())
				assert.Equal(t, tt.expected.IsNull(), s.IsNull())
			}
		})
	}
}

func TestScan(t *testing.T) {
	tests := []struct {
		name         string
		input        any
		expectedVal  string
		expectedNull bool
	}{
		{"scan string", "scanned", "scanned", false},
		{"scan nil", nil, "", true},
		{"scan int", 123, "123", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var s ztype.String
			err := s.Scan(tt.input)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedVal, s.Get())
			assert.Equal(t, tt.expectedNull, s.IsNull())
		})
	}
}

func TestValue(t *testing.T) {
	tests := []struct {
		name         string
		input        ztype.String
		expectedVal  any
		expectedNull bool
	}{
		{"non-null", ztype.NewString("value"), "value", false},
		{"null", ztype.NewNullString(), nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			val, err := tt.input.Value()
			assert.NoError(t, err)
			if tt.expectedVal == nil {
				assert.Nil(t, val)
			} else {
				assert.Equal(t, tt.expectedVal, val)
			}
		})
	}
}

func TestString(t *testing.T) {
	tests := []struct {
		name     string
		input    ztype.String
		expected string
	}{
		{"non-null", ztype.NewString("text"), "text"},
		{"null", ztype.NewNullString(), "<NULL>"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.input.String())
		})
	}
}
