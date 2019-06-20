package astgen

import (
	"fmt"
	"go/ast"
	"go/token"
	"reflect"
)

func buildType(t reflect.Type) (ast.Expr, error) {
	switch t.Kind() {
	case reflect.Bool:
		return &ast.Ident{Name: "bool"}, nil
	case reflect.Int:
		return &ast.Ident{Name: "int"}, nil
	case reflect.Int8:
		return &ast.Ident{Name: "int8"}, nil
	case reflect.Int16:
		return &ast.Ident{Name: "int16"}, nil
	case reflect.Int32:
		return &ast.Ident{Name: "int32"}, nil
	case reflect.Int64:
		return &ast.Ident{Name: "int64"}, nil
	case reflect.Uint:
		return &ast.Ident{Name: "uint"}, nil
	case reflect.Uint8:
		return &ast.Ident{Name: "uint8"}, nil
	case reflect.Uint16:
		return &ast.Ident{Name: "uint16"}, nil
	case reflect.Uint32:
		return &ast.Ident{Name: "uint32"}, nil
	case reflect.Uint64:
		return &ast.Ident{Name: "uint64"}, nil
	case reflect.Float32:
		return &ast.Ident{Name: "float32"}, nil
	case reflect.Float64:
		return &ast.Ident{Name: "float64"}, nil
	case reflect.Complex64:
		return &ast.Ident{Name: "complex64"}, nil
	case reflect.Complex128:
		return &ast.Ident{Name: "complex128"}, nil
	case reflect.String:
		return &ast.Ident{Name: "string"}, nil
	case reflect.Interface:
		return &ast.InterfaceType{Methods: &ast.FieldList{}}, nil
	case reflect.Array:
		elem, err := buildType(t.Elem())
		if err != nil {
			return nil, err
		}
		return &ast.ArrayType{
			Len: &ast.BasicLit{Kind: token.INT, Value: fmt.Sprint(t.Len())},
			Elt: elem,
		}, nil
	case reflect.Slice:
		elem, err := buildType(t.Elem())
		if err != nil {
			return nil, err
		}
		return &ast.ArrayType{Elt: elem}, nil
	case reflect.Map:
		k, err := buildType(t.Key())
		if err != nil {
			return nil, err
		}
		v, err := buildType(t.Elem())
		if err != nil {
			return nil, err
		}
		return &ast.MapType{Key: k, Value: v}, nil
	case reflect.Struct:
		if t.Name() != "" {
			return &ast.Ident{Name: t.Name()}, nil
		}
		fs := make([]*ast.Field, 0, t.NumField())
		var prevType ast.Expr
		for i := 0; i < t.NumField(); i++ {
			sf := t.Field(i)
			t, err := buildType(sf.Type)
			if err != nil {
				return nil, err
			}
			if reflect.DeepEqual(prevType, t) {
				fs[len(fs)-1].Names = append(fs[len(fs)-1].Names, &ast.Ident{Name: sf.Name})
				continue
			}
			fs = append(fs, &ast.Field{
				Names: []*ast.Ident{&ast.Ident{Name: sf.Name}},
				Type:  t,
			})
			prevType = t
		}
		return &ast.StructType{Fields: &ast.FieldList{List: fs}}, nil
	case reflect.Ptr:
		t, err := buildType(t.Elem())
		if err != nil {
			return nil, err
		}
		return &ast.StarExpr{X: t}, nil
	default:
		return nil, &unexpectedTypeError{t}
	}
}
