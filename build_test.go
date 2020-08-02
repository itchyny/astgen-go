package astgen_test

import (
	"bytes"
	"go/parser"
	"go/printer"
	"go/token"
	"math"
	"testing"

	"github.com/itchyny/astgen-go"
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
		src:      [2][1]int{{0}, {1}},
		expected: `[2][1]int{[1]int{0}, [1]int{1}}`,
	},
	{
		name:     "array of array of array",
		src:      [1][1][1]int{{{1}}},
		expected: `[1][1][1]int{[1][1]int{[1]int{1}}}`,
	},
	{
		name:     "slice of int",
		src:      []int{1, 2, 3, 4, 5},
		expected: `[]int{1, 2, 3, 4, 5}`,
	},
	{
		name:     "slice of array of int",
		src:      [][2]int{{1, 2}, {3, 4}},
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
		src:      []map[int]string{{1: "a"}, {2: "b"}},
		expected: `[]map[int]string{map[int]string{1: "a"}, map[int]string{2: "b"}}`,
	},
	{
		name:     "map of int from string",
		src:      map[string]int{"e": 5, "b": 2, "c": 3, "d": 4, "a": 1},
		expected: `map[string]int{"a": 1, "b": 2, "c": 3, "d": 4, "e": 5}`,
	},
	{
		name:     "map of slice of string from int",
		src:      map[int][]string{128: {"Hello", "world!"}, 0: {}},
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
		expected: `(func(x int) *int {
	return &x
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
		expected: `(func(f, b, ba z, x, x2 y) struct {
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
	}{y: 1, z: &f, w: &b, u: &ba, a: &x, b: &x2, c: &x}
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
		expected: `(func(f, b, fo, ba string) map[int]*string {
	return map[int]*string{2: &f, 3: &f, 4: &b, 5: &fo, 7: &ba}
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
		expected: `(func(t, f bool) map[string]*bool {
	return map[string]*bool{"": &t, "a": &t, "c": &f, "foo": &f}
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
		expected: `(func(x int, i int8, i1 int16, i10 int32, i101 int64, u uint, u1 uint8, u10 uint16, u101 uint32, u102 uint64, f float32, x1 float64, c complex64, c1 complex128, in, is interface {
}) map[string]interface {
} {
	return map[string]interface {
	}{"a": interface {
	}(&x), "b": interface {
	}(&i), "c": interface {
	}(&i1), "d": interface {
	}(&i10), "e": interface {
	}(&i101), "f": interface {
	}(&u), "g": interface {
	}(&u1), "h": interface {
	}(&u10), "i": interface {
	}(&u101), "j": interface {
	}(&u102), "k": interface {
	}(&f), "l": interface {
	}(&x1), "m": interface {
	}(&c), "n": interface {
	}(&c1), "o": interface {
	}(&in), "p": interface {
	}(&is)}
})(10, int8(10), int16(10), int32(10), int64(10), uint(10), uint8(10), uint16(10), uint32(10), uint64(10), float32(10), 10.0, complex64((10+0i)), complex128((10+0i)), interface {
}(nil), interface {
}(struct {
}{}))`,
	},
	{
		name: "pointers of pointer",
		src: map[string]interface{}{
			"a": (func(x struct{}) ***struct{} { y := &x; z := &y; return &z })(struct{}{}),
			"b": (func(x bool) **bool { y := &x; return &y })(false),
			"c": (func(x string) ***string { y := &x; z := &y; return &z })(""),
		},
		expected: `(func(s *struct {
}, f bool, x string) map[string]interface {
} {
	s1 := &s
	f1 := &f
	x1 := &x
	x11 := &x1
	return map[string]interface {
	}{"a": interface {
	}(&s1), "b": interface {
	}(&f1), "c": interface {
	}(&x11)}
})(&struct {
}{}, false, "")`,
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
			got, err := astgen.Build(tc.src)
			if err != nil {
				t.Fatalf("should not return error: %s", err)
			}
			buf := new(bytes.Buffer)
			printer.Fprint(buf, token.NewFileSet(), got)
			if buf.String() != tc.expected {
				t.Errorf("expected: %s\ngot: %s", tc.expected, buf.String())
			}
			_, err = parser.ParseExpr(buf.String())
			if err != nil {
				t.Fatalf("should not return error: %s", err)
			}
		})
	}
}
