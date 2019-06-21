package astgen

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/printer"
	"go/token"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

// Build ast from interface{}.
func Build(x interface{}) (ast.Node, error) {
	return (&builder{}).build(reflect.ValueOf(x))
}

type builder struct {
	vars []builderVar
}

type builderVar struct {
	name string
	typ  ast.Expr
	expr ast.Expr
}

func (b *builder) build(v reflect.Value) (ast.Node, error) {
	n, err := b.buildInner(v)
	if err != nil {
		return nil, err
	}
	if len(b.vars) == 0 {
		return n, nil
	}
	t, err := buildType(v.Type())
	if err != nil {
		return nil, err
	}
	params := make([]*ast.Field, len(b.vars))
	args := make([]ast.Expr, len(b.vars))
	for i, bv := range b.vars {
		params[i] = &ast.Field{
			Names: []*ast.Ident{&ast.Ident{Name: bv.name}},
			Type:  bv.typ,
		}
		args[i] = bv.expr
	}
	return &ast.CallExpr{
		Fun: &ast.ParenExpr{
			X: &ast.FuncLit{
				Type: &ast.FuncType{
					Params: &ast.FieldList{List: params},
					Results: &ast.FieldList{
						List: []*ast.Field{
							&ast.Field{Type: t},
						},
					},
				},
				Body: &ast.BlockStmt{
					List: []ast.Stmt{
						&ast.ReturnStmt{
							Results: []ast.Expr{n},
						},
					},
				},
			},
		},
		Args: args,
	}, nil
}

func (b *builder) buildInner(v reflect.Value) (ast.Expr, error) {
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
		if strings.ContainsRune(v.String(), '"') && !strings.ContainsRune(v.String(), '`') {
			s := strings.Replace(v.String(), `"`, "", -1)
			if len(strconv.Quote(s)) == len(s)+2 { // check no escape characters
				return &ast.BasicLit{Kind: token.STRING, Value: "`" + v.String() + "`"}, nil
			}
		}
		return &ast.BasicLit{Kind: token.STRING, Value: strconv.Quote(v.String())}, nil
	case reflect.Interface:
		e, err := b.buildExpr(v.Elem())
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
			w, err := b.buildExpr(v.Index(i))
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
		iter := v.MapRange()
		exprs := make([]ast.Expr, v.Len())
		var i int
		for iter.Next() {
			k, err := b.buildExpr(iter.Key())
			if err != nil {
				return nil, err
			}
			v, err := b.buildExpr(iter.Value())
			if err != nil {
				return nil, err
			}
			exprs[i] = &ast.KeyValueExpr{Key: k, Value: v}
			i++
		}
		sort.Slice(exprs, func(i, j int) bool {
			var buf1, buf2 bytes.Buffer
			printer.Fprint(&buf1, token.NewFileSet(), exprs[i])
			printer.Fprint(&buf2, token.NewFileSet(), exprs[j])
			return buf1.String() < buf2.String()
		})
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
			v, err := b.buildExpr(v.Field(i))
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
		w, err := b.buildExpr(v.Elem())
		if err != nil {
			return nil, err
		}
		if x, ok := w.(*ast.BasicLit); ok {
			t, err := buildType(v.Elem().Type())
			if err != nil {
				return nil, err
			}
			n := b.getVarName(t, x)
			return &ast.UnaryExpr{Op: token.AND, X: &ast.Ident{Name: n}}, nil
		}
		return &ast.UnaryExpr{Op: token.AND, X: w}, nil
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

func (b *builder) buildExpr(v reflect.Value) (ast.Expr, error) {
	w, err := b.buildInner(v)
	if err != nil {
		return nil, err
	}
	e, ok := w.(ast.Expr)
	if !ok {
		return nil, fmt.Errorf("expected ast.Expr but got: %T", w)
	}
	return e, nil
}

func (b *builder) getVarName(t, v ast.Expr) string {
	for _, bv := range b.vars {
		if reflect.DeepEqual(t, bv.typ) && reflect.DeepEqual(v, bv.expr) {
			return bv.name
		}
	}
	name := "x" + strconv.Itoa(len(b.vars))
	bv := builderVar{name: name, typ: t, expr: v}
	b.vars = append(b.vars, bv)
	return name
}
