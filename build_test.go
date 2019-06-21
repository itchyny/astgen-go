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
		name:     "float64",
		src:      3.00,
		expected: `3.0`,
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
		name: "pointer of literal",
		src:  (func(i int) *int { return &i })(42),
		expected: `(func(x4 int) *int {
	return &x4
})(42)`,
	},
	{
		name: "non struct type",
		src: struct {
			x x
			y y
			z *z
			w *z
			u *z
			a *y
			b *y
			c *y
		}{
			y: 1,
			z: (func(s z) *z { return &s })("foo"),
			w: (func(s z) *z { return &s })("bar"),
			u: (func(s z) *z { return &s })("barr"),
			a: (func(i y) *y { return &i })(1),
			b: (func(i y) *y { return &i })(2),
			c: (func(i y) *y { return &i })(1),
		},
		expected: `(func(xf, xb, xba z, x1, x2 y) struct {
	x	x
	y	y
	z, w, u	*z
	a, b, c	*y
} {
	return struct {
		x	x
		y	y
		z, w, u	*z
		a, b, c	*y
	}{y: 1, z: &xf, w: &xb, u: &xba, a: &x1, b: &x2, c: &x1}
})("foo", "bar", "barr", 1, 2)`,
	},
	{
		name: "map of pointers of strings",
		src: map[int]*string{
			3: (func(s string) *string { return &s })("foo"),
			2: (func(s string) *string { return &s })("foo"),
			4: (func(s string) *string { return &s })("bar"),
			5: (func(s string) *string { return &s })("fo"),
			7: (func(s string) *string { return &s })("ba"),
		},
		expected: `(func(xf, xb, xfo, xba string) map[int]*string {
	return map[int]*string{2: &xf, 3: &xf, 4: &xb, 5: &xfo, 7: &xba}
})("foo", "bar", "fo", "ba")`,
	},
	{
		name: "map of pointers of booleans",
		src: map[string]*bool{
			"foo": (func(b bool) *bool { return &b })(false),
			"a":   (func(b bool) *bool { return &b })(true),
			"c":   (func(b bool) *bool { return &b })(false),
			"":    (func(b bool) *bool { return &b })(true),
		},
		expected: `(func(xt, xf bool) map[string]*bool {
	return map[string]*bool{"": &xt, "a": &xt, "c": &xf, "foo": &xf}
})(true, false)`,
	},
	{
		name: "map of pointers of int, uint, float, complex, interface",
		src: map[string]interface{}{
			"a": (func(x int) *int { return &x })(10),
			"b": (func(x int8) *int8 { return &x })(10),
			"c": (func(x int16) *int16 { return &x })(10),
			"d": (func(x int32) *int32 { return &x })(10),
			"e": (func(x int64) *int64 { return &x })(10),
			"f": (func(x uint) *uint { return &x })(10),
			"g": (func(x uint8) *uint8 { return &x })(10),
			"h": (func(x uint16) *uint16 { return &x })(10),
			"i": (func(x uint32) *uint32 { return &x })(10),
			"j": (func(x uint64) *uint64 { return &x })(10),
			"k": (func(x float32) *float32 { return &x })(10),
			"l": (func(x float64) *float64 { return &x })(10),
			"m": (func(x complex64) *complex64 { return &x })(10),
			"n": (func(x complex128) *complex128 { return &x })(10),
			"o": (func(x interface{}) *interface{} { return &x })(nil),
			"p": (func(x interface{}) *interface{} { return &x })(struct{}{}),
		},
		expected: `(func(x1 int, xi int8, xin int16, xint int32, xint1 int64, xu uint, xui uint8, xuin uint16, xuin1 uint32, xuin2 uint64, xf float32, x10 float64, xc complex64, xco complex128, xint2, xint3 interface {
}) map[string]interface {
} {
	return map[string]interface {
	}{"a": interface {
	}(&x1), "b": interface {
	}(&xi), "c": interface {
	}(&xin), "d": interface {
	}(&xint), "e": interface {
	}(&xint1), "f": interface {
	}(&xu), "g": interface {
	}(&xui), "h": interface {
	}(&xuin), "i": interface {
	}(&xuin1), "j": interface {
	}(&xuin2), "k": interface {
	}(&xf), "l": interface {
	}(&x10), "m": interface {
	}(&xc), "n": interface {
	}(&xco), "o": interface {
	}(&xint2), "p": interface {
	}(&xint3)}
})(10, int8(10), int16(10), int32(10), int64(10), uint(10), uint8(10), uint16(10), uint32(10), uint64(10), float32(10), 10.0, complex64((10+0i)), complex128((10+0i)), interface {
}(nil), interface {
}(struct {
}{}))`,
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
