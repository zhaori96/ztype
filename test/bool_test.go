package ztype_test

import (
	"database/sql/driver"
	"encoding/json"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/zhaori96/ztype"
)

func TestBool(t *testing.T) {
	t.Run("Constructors", func(t *testing.T) {
		t.Run("NewBool", func(t *testing.T) {
			b := ztype.NewBool(true)
			require.True(t, b.Get())
			require.False(t, b.IsNull())
		})

		t.Run("NewNullBool", func(t *testing.T) {
			b := ztype.NewNullBool()
			require.True(t, b.IsNull())
		})
	})

	t.Run("ValueManipulation", func(t *testing.T) {
		t.Run("Set", func(t *testing.T) {
			var b ztype.Bool
			b.Set(true)
			require.True(t, b.Get())
			require.False(t, b.IsNull())
		})

		t.Run("SetNull", func(t *testing.T) {
			b := ztype.NewBool(true)
			b.SetNull()
			require.True(t, b.IsNull())
			require.False(t, b.Get())
		})
	})

	t.Run("StateChecks", func(t *testing.T) {
		tests := []struct {
			name     string
			instance ztype.Bool
			isNull   bool
			isZero   bool
		}{
			{"Valid true", ztype.NewBool(true), false, false},
			{"Valid false", ztype.NewBool(false), false, true},
			{"Null", ztype.NewNullBool(), true, true},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				require.Equal(t, tt.isNull, tt.instance.IsNull())
				require.Equal(t, tt.isZero, tt.instance.IsZero())
			})
		}
	})

	t.Run("Serialization", func(t *testing.T) {
		t.Run("MarshalText", func(t *testing.T) {
			tests := []struct {
				name     string
				instance ztype.Bool
				expected []byte
			}{
				{"True", ztype.NewBool(true), []byte("true")},
				{"False", ztype.NewBool(false), []byte("false")},
				{"Null", ztype.NewNullBool(), nil},
			}

			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					data, err := tt.instance.MarshalText()
					require.NoError(t, err)
					require.Equal(t, tt.expected, data)
				})
			}
		})

		t.Run("UnmarshalText", func(t *testing.T) {
			tests := []struct {
				input    string
				expected bool
				valid    bool
			}{
				{"true", true, true},
				{"false", false, true},
				{"invalid", false, false},
			}

			for _, tt := range tests {
				t.Run(tt.input, func(t *testing.T) {
					var b ztype.Bool
					err := b.UnmarshalText([]byte(tt.input))

					if tt.input == "invalid" {
						require.Error(t, err)
						return
					}

					require.NoError(t, err)
					require.Equal(t, tt.expected, b.Get())
					require.Equal(t, tt.valid, !b.IsNull())
					require.True(t, b.Unmarshaled())
				})
			}
		})
	})

	t.Run("JSONHandling", func(t *testing.T) {
		t.Run("MarshalJSON", func(t *testing.T) {
			tests := []struct {
				name     string
				instance ztype.Bool
				expected string
			}{
				{"True", ztype.NewBool(true), "true"},
				{"False", ztype.NewBool(false), "false"},
				{"Null", ztype.NewNullBool(), "null"},
			}

			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					data, err := tt.instance.MarshalJSON()
					require.NoError(t, err)
					require.JSONEq(t, tt.expected, string(data))
				})
			}
		})

		t.Run("UnmarshalJSON", func(t *testing.T) {
			tests := []struct {
				input    string
				expected ztype.Bool
			}{
				{"true", ztype.NewBool(true)},
				{"false", ztype.NewBool(false)},
				{"null", ztype.NewNullBool()},
			}

			for _, tt := range tests {
				t.Run(tt.input, func(t *testing.T) {
					var b ztype.Bool
					err := json.Unmarshal([]byte(tt.input), &b)
					require.NoError(t, err)
					require.True(t, b.Equal(tt.expected))
					require.True(t, b.Unmarshaled())
				})
			}

			t.Run("Invalid", func(t *testing.T) {
				var b ztype.Bool
				err := json.Unmarshal([]byte(`"string"`), &b)
				require.Error(t, err)
			})
		})
	})

	t.Run("DatabaseIntegration", func(t *testing.T) {
		t.Run("Scan", func(t *testing.T) {
			tests := []struct {
				name     string
				input    any
				expected ztype.Bool
			}{
				{"True", true, ztype.NewBool(true)},
				{"False", false, ztype.NewBool(false)},
				{"Null", nil, ztype.NewNullBool()},
			}

			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					var b ztype.Bool
					err := b.Scan(tt.input)
					require.NoError(t, err)
					require.True(t, b.Equal(tt.expected))
				})
			}

			t.Run("InvalidType", func(t *testing.T) {
				var b ztype.Bool
				err := b.Scan("string")
				require.Error(t, err)
			})
		})

		t.Run("Value", func(t *testing.T) {
			tests := []struct {
				name     string
				instance ztype.Bool
				expected driver.Value
			}{
				{"True", ztype.NewBool(true), true},
				{"False", ztype.NewBool(false), false},
				{"Null", ztype.NewNullBool(), nil},
			}

			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					val, err := tt.instance.Value()
					require.NoError(t, err)
					require.Equal(t, tt.expected, val)
				})
			}
		})
	})

	t.Run("Comparisons", func(t *testing.T) {
		t.Run("Equal", func(t *testing.T) {
			tests := []struct {
				a        ztype.Bool
				b        ztype.Bool
				expected bool
			}{
				{ztype.NewBool(true), ztype.NewBool(true), true},
				{ztype.NewBool(true), ztype.NewBool(false), false},
				{ztype.NewNullBool(), ztype.NewNullBool(), true},
				{ztype.NewBool(true), ztype.NewNullBool(), false},
			}

			for i, tt := range tests {
				t.Run(strconv.Itoa(i), func(t *testing.T) {
					require.Equal(t, tt.expected, tt.a.Equal(tt.b))
				})
			}
		})

		t.Run("EqualRaw", func(t *testing.T) {
			tests := []struct {
				instance ztype.Bool
				input    bool
				expected bool
			}{
				{ztype.NewBool(true), true, true},
				{ztype.NewNullBool(), true, false},
				{ztype.NewNullBool(), false, true},
			}

			for i, tt := range tests {
				t.Run(strconv.Itoa(i), func(t *testing.T) {
					require.Equal(t, tt.expected, tt.instance.EqualRaw(tt.input))
				})
			}
		})
	})

	t.Run("StringRepresentation", func(t *testing.T) {
		tests := []struct {
			instance ztype.Bool
			expected string
		}{
			{ztype.NewBool(true), "true"},
			{ztype.NewBool(false), "false"},
			{ztype.NewNullBool(), "<NULL>"},
		}

		for _, tt := range tests {
			t.Run(tt.expected, func(t *testing.T) {
				require.Equal(t, tt.expected, tt.instance.String())
			})
		}
	})

	t.Run("UnmarshaledTracking", func(t *testing.T) {
		t.Run("DefaultState", func(t *testing.T) {
			b := ztype.NewBool(true)
			require.False(t, b.Unmarshaled())
		})

		t.Run("AfterUnmarshal", func(t *testing.T) {
			var b ztype.Bool
			json.Unmarshal([]byte("true"), &b)
			require.True(t, b.Unmarshaled())
		})

		t.Run("ManualSet", func(t *testing.T) {
			var b ztype.Bool
			b.SetUnmarshaled(true)
			require.True(t, b.Unmarshaled())
		})
	})
}
