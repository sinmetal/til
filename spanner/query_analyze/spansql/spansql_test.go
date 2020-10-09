package spansql

import (
	"testing"

	"cloud.google.com/go/spanner/spansql"
	"github.com/k0kubun/pp"
)

func TestWalk2(t *testing.T) {
	query, err := spansql.ParseQuery(`SELECT id, createdAt FROM customers WHERE ID = "sinmetal" AND LAST = 1 LIMIT 3`)
	if err != nil {
		t.Fatal(err)
	}

	where := query.Select.Where
	var v spansql.LogicalOp = where.(spansql.LogicalOp)
	pp.Println(v)
}

func TestWalk(t *testing.T) {
	query, err := spansql.ParseQuery(`SELECT id, createdAt FROM customers WHERE ID = "sinmetal" AND LAST = 1 LIMIT 3`)
	if err != nil {
		t.Fatal(err)
	}
	query.Limit = HAMMER_INTEGER_LITERAL

	newQuery := Walk(query, func(node interface{}) interface{} {
		switch v := node.(type) {
		case spansql.StringLiteral:
			return HAMMER_STRING_LITERAL
		case spansql.IntegerLiteral:
			return HAMMER_INTEGER_LITERAL
		default:
			return v
		}
	}).(spansql.Query)

	t.Log(newQuery.SQL())
	pp.Println(newQuery)
}
