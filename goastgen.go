package goastgen

import (
	"fmt"
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
	case reflect.Bool:
		if v.Bool() {
			return &ast.Ident{Name: "true"}
		}
		return &ast.Ident{Name: "false"}
	case reflect.Int:
		return &ast.BasicLit{Kind: token.INT, Value: fmt.Sprint(v.Int())}
	case reflect.Float64:
		return &ast.BasicLit{Kind: token.FLOAT, Value: fmt.Sprint(v.Float())}
	case reflect.String:
		return &ast.BasicLit{Kind: token.STRING, Value: strconv.Quote(v.String())}
	default:
		return nil
	}
}
