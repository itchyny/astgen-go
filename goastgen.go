package goastgen

import "go/ast"

// Build ast from interface{}.
func Build(x interface{}) ast.Node {
	if x == nil {
		return &ast.Ident{Name: "nil"}
	}
	return nil
}
