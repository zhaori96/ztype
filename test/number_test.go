package ztype_test

import (
	"encoding/json"
	"math"
	"reflect"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/zhaori96/ztype"
)

type numericTestCase struct {
	name        string
	typ         string
	validVal    any
	otherVal    any
	zeroVal     any
	minVal      any
	maxVal      any
	floatInput  string
	intInput    string
	overflowAdd any
}

var numericTestCases = []numericTestCase{
	{
		name:        "int",
		typ:         "int",
		validVal:    42,
		otherVal:    100,
		zeroVal:     0,
		minVal:      math.MinInt32,
		maxVal:      math.MaxInt32,
		floatInput:  "123.5",
		intInput:    "123",
		overflowAdd: 142,
	},
	{
		name:        "int8",
		typ:         "int8",
		validVal:    int8(20),
		otherVal:    int8(42),
		zeroVal:     int8(0),
		minVal:      math.MinInt8,
		maxVal:      math.MaxInt8,
		floatInput:  "64.5",
		intInput:    "64",
		overflowAdd: int8(62),
	},
	{
		name:        "uint",
		typ:         "uint",
		validVal:    uint(42),
		otherVal:    uint(100),
		zeroVal:     uint(0),
		minVal:      uint(0),
		maxVal:      uint(math.MaxUint32),
		floatInput:  "255.5",
		intInput:    "255",
		overflowAdd: uint(142),
	},
	{
		name:       "float32",
		typ:        "float32",
		validVal:   float32(3.14),
		otherVal:   float32(6.28),
		zeroVal:    float32(0),
		minVal:     -math.MaxFloat32,
		maxVal:     math.MaxFloat32,
		floatInput: "123.456",
		intInput:   "123",
	},
	{
		name:       "float64",
		typ:        "float64",
		validVal:   3.141592653589793,
		otherVal:   6.283185307179586,
		zeroVal:    0.0,
		minVal:     -math.MaxFloat64,
		maxVal:     math.MaxFloat64,
		floatInput: "123.456789",
		intInput:   "123",
	},
}

func TestNumericConstructors(t *testing.T) {
	for _, tc := range numericTestCases {
		t.Run(tc.typ, func(t *testing.T) {
			switch tc.typ {
			case "int":
				testConstructor[int](t, tc)
			case "int8":
				testConstructor[int8](t, tc)
			case "uint":
				testConstructor[uint](t, tc)
			case "float32":
				testConstructor[float32](t, tc)
			case "float64":
				testConstructor[float64](t, tc)
			}
		})
	}
}

func testConstructor[T ztype.NumberType](t *testing.T, tc numericTestCase) {
	t.Run("NewNumber", func(t *testing.T) {
		n := ztype.NewNumber(tc.validVal.(T))
		assert.Equal(t, tc.validVal, n.Get())
		assert.False(t, n.IsNull())
	})

	t.Run("NewNullNumber", func(t *testing.T) {
		n := ztype.NewNullNumber[T]()
		assert.Equal(t, tc.zeroVal, n.Get())
		assert.True(t, n.IsNull())
	})
}

func TestNumericSetters(t *testing.T) {
	for _, tc := range numericTestCases {
		t.Run(tc.typ, func(t *testing.T) {
			switch tc.typ {
			case "int":
				testSetters[int](t, tc)
			case "int8":
				testSetters[int8](t, tc)
			case "uint":
				testSetters[uint](t, tc)
			case "float32":
				testSetters[float32](t, tc)
			case "float64":
				testSetters[float64](t, tc)
			}
		})
	}
}

func testSetters[T ztype.NumberType](t *testing.T, tc numericTestCase) {
	t.Run("Set", func(t *testing.T) {
		var n ztype.Numeric[T]
		n.Set(tc.validVal.(T))
		assert.Equal(t, tc.validVal, n.Get())
		assert.False(t, n.IsNull())
	})

	t.Run("SetNull", func(t *testing.T) {
		n := ztype.NewNumber(tc.validVal.(T))
		n.SetNull()
		assert.Equal(t, tc.zeroVal, n.Get())
		assert.True(t, n.IsNull())
	})
}

func TestNumericOperations(t *testing.T) {
	for _, tc := range numericTestCases {
		t.Run(tc.typ, func(t *testing.T) {
			switch tc.typ {
			case "int":
				testOperations[int](t, tc)
			case "int8":
				testOperations[int8](t, tc)
			case "uint":
				testOperations[uint](t, tc)
			case "float32":
				testOperations[float32](t, tc)
			case "float64":
				testOperations[float64](t, tc)
			}
		})
	}
}

func testOperations[T ztype.NumberType](t *testing.T, tc numericTestCase) {
	valid := ztype.NewNumber(tc.validVal.(T))
	other := ztype.NewNumber(tc.otherVal.(T))
	null := ztype.NewNullNumber[T]()

	t.Run("Add", func(t *testing.T) {
		result := valid.Add(other)
		if tc.overflowAdd != nil {
			assert.Equal(t, tc.overflowAdd, result.Get())
		} else {
			assert.Equal(t, add(tc.validVal, tc.otherVal), result.Get())
		}
		assert.True(t, valid.Add(null).IsNull())
	})

	t.Run("Compare", func(t *testing.T) {
		result, err := valid.Compare(other)
		assert.NoError(t, err)
		assert.Equal(t, -1, result)

		_, err = valid.Compare(null)
		assert.Error(t, err)
	})

	t.Run("Division", func(t *testing.T) {
		_, err := valid.SafeDiv(ztype.NewNumber(tc.zeroVal.(T)))
		assert.Error(t, err)
	})
}

func add(a, b any) any {
	switch a := a.(type) {
	case int:
		return a + b.(int)
	case int8:
		return a + b.(int8)
	case uint:
		return a + b.(uint)
	case float32:
		return a + b.(float32)
	case float64:
		return a + b.(float64)
	}
	return nil
}

func TestNumericJSON(t *testing.T) {
	for _, tc := range numericTestCases {
		t.Run(tc.typ, func(t *testing.T) {
			switch tc.typ {
			case "int":
				testJSON[int](t, tc)
			case "int8":
				testJSON[int8](t, tc)
			case "uint":
				testJSON[uint](t, tc)
			case "float32":
				testJSON[float32](t, tc)
			case "float64":
				testJSON[float64](t, tc)
			}
		})
	}
}

func testJSON[T ztype.NumberType](t *testing.T, tc numericTestCase) {
	t.Run("Marshal/Unmarshal", func(t *testing.T) {
		n := ztype.NewNumber(tc.validVal.(T))
		data, err := n.MarshalJSON()
		assert.NoError(t, err)

		var unmarshaled ztype.Numeric[T]
		assert.NoError(t, json.Unmarshal(data, &unmarshaled))
		assert.Equal(t, n.Get(), unmarshaled.Get())
	})

	t.Run("FloatInput", func(t *testing.T) {
		var n ztype.Numeric[T]
		err := json.Unmarshal([]byte(tc.floatInput), &n)

		switch any(tc.zeroVal).(type) {
		case int, int8, uint:
			assert.Error(t, err, "Should reject non-integer float")
		case float32, float64:
			assert.NoError(t, err)
			assert.InEpsilon(t, parseFloat(tc.floatInput), n.Get(), 1e-6)
		}
	})

	t.Run("IntInput", func(t *testing.T) {
		var n ztype.Numeric[T]
		err := json.Unmarshal([]byte(tc.intInput), &n)
		assert.NoError(t, err)
		assert.Equal(t, T(parseInt(tc.intInput)), n.Get())
	})
}

func parseFloat(s string) float64 {
	f, _ := strconv.ParseFloat(s, 64)
	return f
}

func parseInt(s string) int64 {
	i, _ := strconv.ParseInt(s, 10, 64)
	return i
}

func TestDatabaseIntegration(t *testing.T) {
	for _, tc := range numericTestCases {
		t.Run(tc.typ, func(t *testing.T) {
			switch tc.typ {
			case "int":
				testDatabase[int](t, tc)
			case "int8":
				testDatabase[int8](t, tc)
			case "uint":
				testDatabase[uint](t, tc)
			case "float32":
				testDatabase[float32](t, tc)
			case "float64":
				testDatabase[float64](t, tc)
			}
		})
	}
}

func testDatabase[T ztype.NumberType](t *testing.T, tc numericTestCase) {
	t.Run("Scan/Value", func(t *testing.T) {
		var n ztype.Numeric[T]
		v := reflect.ValueOf(tc.validVal)

		switch v.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			err := n.Scan(v.Int())
			assert.NoError(t, err)

		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			err := n.Scan(int64(v.Uint()))
			assert.NoError(t, err)

		case reflect.Float32, reflect.Float64:
			err := n.Scan(v.Float())
			assert.NoError(t, err)

		default:
			t.Fatalf("Tipo não suportado: %s", v.Kind())
		}
		val, err := n.Value()
		assert.NoError(t, err)
		var expected interface{}
		switch v.Kind() {
		case reflect.Int, reflect.Int8:
			expected = int64(v.Int())
		case reflect.Uint:
			expected = int64(v.Uint())
		case reflect.Float32:
			expected = float64(v.Float())
		default:
			expected = tc.validVal
		}

		assert.Equal(t, expected, val)
	})
}
