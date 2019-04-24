package goastgen

import (
	"go/ast"
	"go/token"
	"reflect"
	"strconv"
)

// Build ast from interface{}.
func Build(x interface{}) ast.Node {
	if x == nil {
		return &ast.Ident{Name: "nil"}
	}
	v := reflect.ValueOf(x)
	switch v.Kind() {
	case reflect.String:
		return &ast.BasicLit{Kind: token.STRING, Value: strconv.Quote(v.String())}
	default:
		return nil
	}
}
