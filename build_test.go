package astgen

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
	{
		name:     "string containing double quote",
		src:      `"hello", "こんにちは"`,
		expected: "`\"hello\", \"こんにちは\"`",
	},
	{
		name:     "int array",
		src:      [3]int{-128, 0, 128},
		expected: `[3]int{-128, 0, 128}`,
	},
	{
		name:     "string array",
		src:      [2]string{"Hello", "world!"},
		expected: `[2]string{"Hello", "world!"}`,
	},
	{
		name:     "array of array",
		src:      [2][1]int{[1]int{0}, [1]int{1}},
		expected: `[2][1]int{[1]int{0}, [1]int{1}}`,
	},
	{
		name:     "array of array of array",
		src:      [1][1][1]int{[1][1]int{[1]int{1}}},
		expected: `[1][1][1]int{[1][1]int{[1]int{1}}}`,
	},
	{
		name:     "slice of int",
		src:      []int{1, 2, 3, 4, 5},
		expected: `[]int{1, 2, 3, 4, 5}`,
	},
	{
		name:     "slice of array of int",
		src:      [][2]int{[2]int{1, 2}, [2]int{3, 4}},
		expected: `[][2]int{[2]int{1, 2}, [2]int{3, 4}}`,
	},
	{
		name: "slice of interface",
		src:  []interface{}{1, "a", nil, false, true},
		expected: `[]interface {
}{interface {
}(1), interface {
}("a"), interface {
}(nil), interface {
}(false), interface {
}(true)}`,
	},
	{
		name:     "slice of map",
		src:      []map[int]string{map[int]string{1: "a"}, map[int]string{2: "b"}},
		expected: `[]map[int]string{map[int]string{1: "a"}, map[int]string{2: "b"}}`,
	},
	{
		name:     "map of int from string",
		src:      map[string]int{"e": 5, "b": 2, "c": 3, "d": 4, "a": 1},
		expected: `map[string]int{"a": 1, "b": 2, "c": 3, "d": 4, "e": 5}`,
	},
	{
		name:     "map of slice of string from int",
		src:      map[int][]string{128: []string{"Hello", "world!"}, 0: []string{}},
		expected: `map[int][]string{0: []string{}, 128: []string{"Hello", "world!"}}`,
	},
	{
		name: "map of interface from string",
		src:  map[string]interface{}{"abcde": 128, "42": []interface{}{}},
		expected: `map[string]interface {
}{"42": interface {
}([]interface {
}{}), "abcde": interface {
}(128)}`,
	},
	{
		name: "empty struct",
		src:  struct{}{},
		expected: `struct {
}{}`,
	},
	{
		name: "struct",
		src: struct {
			foo, bar int
			baz, qux string
			s        []interface{}
			m        map[int]interface{}
		}{foo: 1, baz: "bar", m: map[int]interface{}{1: 128}},
		expected: `struct {
	foo, bar	int
	baz, qux	string
	s		[]interface {
	}
	m	map[int]interface {
	}
}{foo: 1, baz: "bar", m: map[int]interface {
}{1: interface {
}(128)}}`,
	},
	{
		name:     "struct pointer",
		src:      &x{name: "foo"},
		expected: `&x{name: "foo"}`,
	},
	{
		name: "nameless struct pointer",
		src: &struct {
			name string
			ptr1 *int
			ptr2 *int `x:"t,omitempty"`
		}{name: "foo"},
		expected: `&struct {
	name	string
	ptr1	*int
	ptr2	*int	` + "`" + `x:"t,omitempty"` + "`" + `
}{name: "foo"}`,
	},
	{
		name: "non struct type",
		src: struct {
			x x
			y y
			z *z
			w *z
			a *y
			b *y
			c *y
		}{
			y: 1,
			z: (func(s z) *z { return &s })("foo"),
			w: (func(s z) *z { return &s })("bar"),
			a: (func(i y) *y { return &i })(1),
			b: (func(i y) *y { return &i })(2),
			c: (func(i y) *y { return &i })(1),
		},
		expected: `(func(x0 z, x1 z, x2 y, x3 y) struct {
	x	x
	y	y
	z, w	*z
	a, b, c	*y
} {
	return struct {
		x	x
		y	y
		z, w	*z
		a, b, c	*y
	}{y: 1, z: &x0, w: &x1, a: &x2, b: &x3, c: &x2}
})("foo", "bar", 1, 2)`,
	},
}

type x struct {
	name string
	ptr  *int
}

type y int
type z string

func TestBuild(t *testing.T) {
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := Build(tc.src)
			assert.NoError(t, err)
			buf := new(bytes.Buffer)
			printer.Fprint(buf, token.NewFileSet(), got)
			assert.Equal(t, tc.expected, buf.String())
			_, err = parser.ParseExpr(buf.String())
			assert.NoError(t, err)
		})
	}
}
