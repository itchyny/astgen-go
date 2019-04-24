package goastgen

import (
	"bytes"
	"go/parser"
	"go/printer"
	"go/token"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testCases = []struct {
	name     string
	src      interface{}
	expected string
}{
	{
		name:     "nil",
		src:      nil,
		expected: `nil`,
	},
	{
		name:     "true",
		src:      true,
		expected: `true`,
	},
	{
		name:     "false",
		src:      false,
		expected: `false`,
	},
	{
		name:     "int",
		src:      16777216,
		expected: `16777216`,
	},
	{
		name:     "int8",
		src:      int8(math.MinInt8),
		expected: `int8(-128)`,
	},
	{
		name:     "int16",
		src:      int16(math.MaxInt16),
		expected: `int16(32767)`,
	},
	{
		name:     "int32",
		src:      int32(math.MinInt32),
		expected: `int32(-2147483648)`,
	},
	{
		name:     "int64",
		src:      int64(math.MaxInt64),
		expected: `int64(9223372036854775807)`,
	},
	{
		name:     "uint",
		src:      uint(0xffffffffffffffff),
		expected: `uint(18446744073709551615)`,
	},
	{
		name:     "uint8",
		src:      uint8(math.MaxUint8),
		expected: `uint8(255)`,
	},
	{
		name:     "uint16",
		src:      uint16(math.MaxUint16),
		expected: `uint16(65535)`,
	},
	{
		name:     "uint32",
		src:      uint32(math.MaxUint32),
		expected: `uint32(4294967295)`,
	},
	{
		name:     "uint64",
		src:      uint64(math.MaxUint64),
		expected: `uint64(18446744073709551615)`,
	},
	{
		name:     "float32",
		src:      float32(3.125),
		expected: `float32(3.125)`,
	},
	{
		name:     "float64",
		src:      3.14156,
		expected: `3.14156`,
	},
	{
		name:     "complex64",
		src:      complex64(1 - 2i),
		expected: `complex64((1-2i))`,
	},
	{
		name:     "complex128",
		src:      -3.14156 + 2.71828i,
		expected: `complex128((-3.14156+2.71828i))`,
	},
	{
		name:     "string",
		src:      "Hello, world!",
		expected: `"Hello, world!"`,
	},
	{
		name: "string with new lines",
		src: `こんにちは
					世界
					☆ミ`,
		expected: `"こんにちは\n\t\t\t\t\t世界\n\t\t\t\t\t☆ミ"`,
	},
}

func TestBuild(t *testing.T) {
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := Build(tc.src)
			buf := new(bytes.Buffer)
			printer.Fprint(buf, token.NewFileSet(), got)
			assert.Equal(t, tc.expected, buf.String())
			_, err = parser.ParseExpr(buf.String())
			assert.NoError(t, err)
		})
	}
}
