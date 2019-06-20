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
	return build(reflect.ValueOf(x))
}

func build(v reflect.Value) (ast.Node, error) {
	switch v.Kind() {
	case reflect.Invalid:
		return &ast.Ident{Name: "nil"}, nil
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
	case reflect.Interface:
		e, err := buildExpr(v.Elem())
		if err != nil {
			return nil, err
		}
		t, err := buildType(v.Type())
		if err != nil {
			return nil, err
		}
		return &ast.CallExpr{Fun: t, Args: []ast.Expr{e}}, nil
	case reflect.Array, reflect.Slice:
		exprs := make([]ast.Expr, v.Len())
		for i := 0; i < v.Len(); i++ {
			w, err := buildExpr(v.Index(i))
			if err != nil {
				return nil, err
			}
			exprs[i] = w
		}
		t, err := buildType(v.Type())
		if err != nil {
			return nil, err
		}
		return &ast.CompositeLit{Type: t, Elts: exprs}, nil
	case reflect.Map:
		exprs := make([]ast.Expr, v.Len())
		iter := v.MapRange()
		var i int
		for iter.Next() {
			k, err := buildExpr(iter.Key())
			if err != nil {
				return nil, err
			}
			v, err := buildExpr(iter.Value())
			if err != nil {
				return nil, err
			}
			exprs[i] = &ast.KeyValueExpr{Key: k, Value: v}
			i++
		}
		t, err := buildType(v.Type())
		if err != nil {
			return nil, err
		}
		return &ast.CompositeLit{Type: t, Elts: exprs}, nil
	case reflect.Struct:
		exprs := make([]ast.Expr, 0, v.NumField())
		for i := 0; i < v.NumField(); i++ {
			if isZero(v.Field(i)) {
				continue
			}
			k := &ast.Ident{Name: v.Type().Field(i).Name}
			v, err := buildExpr(v.Field(i))
			if err != nil {
				return nil, err
			}
			exprs = append(exprs, &ast.KeyValueExpr{Key: k, Value: v})
		}
		t, err := buildType(v.Type())
		if err != nil {
			return nil, err
		}
		return &ast.CompositeLit{Type: t, Elts: exprs}, nil
	case reflect.Ptr:
		v, err := buildExpr(v.Elem())
		if err != nil {
			return nil, err
		}
		return &ast.UnaryExpr{Op: token.AND, X: v}, nil
	default:
		return nil, &unexpectedTypeError{v.Type()}
	}
}

type unexpectedTypeError struct{ t reflect.Type }

func (err *unexpectedTypeError) Error() string {
	return fmt.Sprintf("unexpected type: %s", err.t.Kind())
}

func callExpr(kind token.Token, name, value string) *ast.CallExpr {
	return &ast.CallExpr{
		Fun: &ast.Ident{Name: name},
		Args: []ast.Expr{
			&ast.BasicLit{Kind: kind, Value: value},
		},
	}
}

func buildExpr(v reflect.Value) (ast.Expr, error) {
	w, err := build(v)
	if err != nil {
		return nil, err
	}
	e, ok := w.(ast.Expr)
	if !ok {
		return nil, fmt.Errorf("expected ast.Expr but got: %T", w)
	}
	return e, nil
}
