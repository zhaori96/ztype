package ztype_test

import (
	"database/sql/driver"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/zhaori96/ztype"
)

func TestByte(t *testing.T) {
	t.Run("Constructor", func(t *testing.T) {
		tests := []struct {
			name     string
			creator  func() ztype.Byte
			expected byte
			isNull   bool
		}{
			{
				name: "NewByte valid",
				creator: func() ztype.Byte {
					return ztype.NewByte(100)
				},
				expected: 100,
				isNull:   false,
			},
			{
				name: "NewNullByte",
				creator: func() ztype.Byte {
					return ztype.NewNullByte()
				},
				expected: 0,
				isNull:   true,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				b := tt.creator()
				require.Equal(t, tt.expected, b.Get())
				require.Equal(t, tt.isNull, b.IsNull())
			})
		}
	})

	t.Run("ValueSemantics", func(t *testing.T) {
		t.Run("Set", func(t *testing.T) {
			var b ztype.Byte
			b.Set(200)
			require.Equal(t, byte(200), b.Get())
			require.False(t, b.IsNull())
		})

		t.Run("SetNull", func(t *testing.T) {
			var b ztype.Byte
			b.Set(50)
			b.SetNull()
			require.True(t, b.IsNull())
			require.Equal(t, byte(0), b.Get())
		})
	})

	t.Run("JSON", func(t *testing.T) {
		tests := []struct {
			name        string
			jsonInput   string
			expected    ztype.Byte
			expectError bool
		}{
			{
				name:        "Valid number",
				jsonInput:   `123`,
				expected:    ztype.NewByte(123),
				expectError: false,
			},
			{
				name:        "Explicit null",
				jsonInput:   `null`,
				expected:    ztype.NewNullByte(),
				expectError: false,
			},
			{
				name:        "Invalid type",
				jsonInput:   `"abc"`,
				expected:    ztype.NewNullByte(),
				expectError: true,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				var b ztype.Byte
				err := json.Unmarshal([]byte(tt.jsonInput), &b)

				if tt.expectError {
					require.Error(t, err)
				} else {
					require.NoError(t, err)
					require.True(t, b.Unmarshaled())
					require.True(t, b.Equal(tt.expected))
				}
			})
		}

		t.Run("Marshal", func(t *testing.T) {
			tests := []struct {
				name     string
				input    ztype.Byte
				expected string
			}{
				{
					name:     "Valid value",
					input:    ztype.NewByte(255),
					expected: "255",
				},
				{
					name:     "Null value",
					input:    ztype.NewNullByte(),
					expected: "null",
				},
			}

			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					result, err := tt.input.MarshalJSON()
					require.NoError(t, err)
					require.JSONEq(t, tt.expected, string(result))
				})
			}
		})
	})

	t.Run("Database", func(t *testing.T) {
		t.Run("Scan", func(t *testing.T) {
			tests := []struct {
				name        string
				input       any
				expected    ztype.Byte
				expectError bool
			}{
				{
					name:        "Valid int64",
					input:       int64(150),
					expected:    ztype.NewByte(150),
					expectError: false,
				},
				{
					name:        "Null input",
					input:       nil,
					expected:    ztype.NewNullByte(),
					expectError: false,
				},
				{
					name:        "Invalid type",
					input:       "invalid",
					expected:    ztype.NewNullByte(),
					expectError: true,
				},
			}

			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					var b ztype.Byte
					err := b.Scan(tt.input)

					if tt.expectError {
						require.Error(t, err)
					} else {
						require.NoError(t, err)
						require.True(t, b.Equal(tt.expected))
					}
				})
			}
		})

		t.Run("Value", func(t *testing.T) {
			tests := []struct {
				name     string
				input    ztype.Byte
				expected driver.Value
			}{
				{
					name:     "Valid value",
					input:    ztype.NewByte(200),
					expected: int64(200),
				},
				{
					name:     "Null value",
					input:    ztype.NewNullByte(),
					expected: nil,
				},
			}

			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					result, err := tt.input.Value()
					require.NoError(t, err)
					require.Equal(t, tt.expected, result)
				})
			}
		})
	})

	t.Run("EdgeCases", func(t *testing.T) {
		t.Run("OverflowProtection", func(t *testing.T) {
			var b ztype.Byte
			err := b.UnmarshalJSON([]byte("256"))
			require.Error(t, err)
			require.True(t, b.IsNull())
		})

		t.Run("ZeroUnmarshaled", func(t *testing.T) {
			var b ztype.Byte
			require.False(t, b.Unmarshaled())
			b.SetUnmarshaled(true)
			require.True(t, b.Unmarshaled())
		})
	})

	t.Run("StringRepresentation", func(t *testing.T) {
		tests := []struct {
			name     string
			input    ztype.Byte
			expected string
		}{
			{
				name:     "Valid value",
				input:    ztype.NewByte(42),
				expected: "42",
			},
			{
				name:     "Null value",
				input:    ztype.NewNullByte(),
				expected: "<NULL>",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				require.Equal(t, tt.expected, tt.input.String())
			})
		}
	})
}
