package goastgen

import (
	"bytes"
	"go/parser"
	"go/printer"
	"go/token"
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
