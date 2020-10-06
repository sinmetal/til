package memefish

import (
	"testing"

	"github.com/MakeNowJust/memefish/pkg/ast"
	"github.com/MakeNowJust/memefish/pkg/parser"
	"github.com/MakeNowJust/memefish/pkg/token"
	"github.com/k0kubun/pp"
)

func TestParseUnParse(t *testing.T) {
	// Create a new Parser instance.
	file := &token.File{
		Buffer: `SELECT * FROM customers WHERE ID = "sinmetal" AND LAST = 1`,
	}
	p := &parser.Parser{
		Lexer: &parser.Lexer{File: file},
	}

	// Do parsing!
	// func Walk(node ast.Node, f func(ast.Node) error) error
	stmt, err := p.ParseQuery()
	if err != nil {
		t.Fatal(err)
	}
	switch v := stmt.Query.(type) {
	case *ast.Select:
		switch v2 := v.Where.Expr.(type) {
		case *ast.BinaryExpr:
			switch left := v2.Left.(type) {
			case *ast.BinaryExpr:
			case *ast.StringLiteral:
				left.Value = "XXXXX"
			case *ast.IntLiteral:
				left.Value = "0"
			default:
			}
		}
	default:
	}

	// Show AST.
	t.Log("AST")
	_, _ = pp.Println(stmt)

	// Unparse AST to SQL source string.
	t.Log("Unparse")
	t.Log(stmt.SQL())
}
