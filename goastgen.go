package goastgen

import (
	"fmt"
	"go/ast"
	"go/token"
	"reflect"
	"strconv"
)

// Build ast from interface{}.
func Build(x interface{}) (ast.Node, error) {
	if x == nil {
		return &ast.Ident{Name: "nil"}, nil
	}
	v := reflect.ValueOf(x)
	return build(v)
}

func build(v reflect.Value) (ast.Node, error) {
	switch v.Kind() {
	case reflect.Bool:
		if v.Bool() {
			return &ast.Ident{Name: "true"}, nil
		}
		return &ast.Ident{Name: "false"}, nil
	case reflect.Int:
		return &ast.BasicLit{Kind: token.INT, Value: fmt.Sprint(v.Int())}, nil
	case reflect.Int8:
		return callExpr(token.INT, "int8", fmt.Sprint(v.Int())), nil
	case reflect.Int16:
		return callExpr(token.INT, "int16", fmt.Sprint(v.Int())), nil
	case reflect.Int32:
		return callExpr(token.INT, "int32", fmt.Sprint(v.Int())), nil
	case reflect.Int64:
		return callExpr(token.INT, "int64", fmt.Sprint(v.Int())), nil
	case reflect.Uint:
		return callExpr(token.INT, "uint", fmt.Sprint(v.Uint())), nil
	case reflect.Uint8:
		return callExpr(token.INT, "uint8", fmt.Sprint(v.Uint())), nil
	case reflect.Uint16:
		return callExpr(token.INT, "uint16", fmt.Sprint(v.Uint())), nil
	case reflect.Uint32:
		return callExpr(token.INT, "uint32", fmt.Sprint(v.Uint())), nil
	case reflect.Uint64:
		return callExpr(token.INT, "uint64", fmt.Sprint(v.Uint())), nil
	case reflect.Float32:
		return callExpr(token.FLOAT, "float32", fmt.Sprint(v.Float())), nil
	case reflect.Float64:
		return &ast.BasicLit{Kind: token.FLOAT, Value: fmt.Sprint(v.Float())}, nil
	case reflect.Complex64:
		return callExpr(token.FLOAT, "complex64", fmt.Sprint(v.Complex())), nil
	case reflect.Complex128:
		return callExpr(token.FLOAT, "complex128", fmt.Sprint(v.Complex())), nil
	case reflect.String:
		return &ast.BasicLit{Kind: token.STRING, Value: strconv.Quote(v.String())}, nil
	case reflect.Array:
		exprs := make([]ast.Expr, v.Len())
		for i := 0; i < v.Len(); i++ {
			v, err := build(v.Index(i))
			if err != nil {
				return nil, err
			}
			w, ok := v.(ast.Expr)
			if !ok {
				return nil, fmt.Errorf("expected ast.Expr but got: %T", v)
			}
			exprs[i] = w
		}
		return &ast.CompositeLit{
			Type: &ast.ArrayType{
				Len: &ast.BasicLit{Kind: token.INT, Value: fmt.Sprint(v.Len())},
				Elt: &ast.Ident{Name: v.Type().Elem().Name()},
			},
			Elts: exprs,
		}, nil
	default:
		return nil, unexpectedTypeError(v)
	}
}

type unexpectedTypeError reflect.Value

func (err unexpectedTypeError) Error() string {
	return fmt.Sprintf("unexpected type: %s", reflect.Value(err).Kind())
}

func callExpr(kind token.Token, name, value string) *ast.CallExpr {
	return &ast.CallExpr{
		Fun: &ast.Ident{Name: name},
		Args: []ast.Expr{
			&ast.BasicLit{Kind: kind, Value: value},
		},
	}
}
