package memefish

import (
	"fmt"
	"testing"

	"github.com/MakeNowJust/memefish/pkg/ast"
	"github.com/MakeNowJust/memefish/pkg/parser"
	"github.com/MakeNowJust/memefish/pkg/token"
)

func TestWalk(t *testing.T) {
	// Create a new Parser instance.
	file := &token.File{
		Buffer: `SELECT 1 FROM Hoge WHERE column = "hello"`,
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
	if err := Walk(stmt, func(node ast.Node) error {
		switch v := node.(type) {
		case *ast.StringLiteral:
			v.Value = "XXXXXX"
			fmt.Printf("change %T\n", v)
		default:
			fmt.Printf("通過 %T\n", v)
		}
		return nil
	}); err != nil {
		t.Fatal(err)
	}

	t.Log("---AFTER---")
	t.Log(stmt.SQL())
}
