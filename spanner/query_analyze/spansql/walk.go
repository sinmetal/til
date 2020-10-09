package spansql

import (
	"cloud.google.com/go/spanner/spansql"
)

const (
	HAMMER_STRING_LITERAL  spansql.StringLiteral  = "HAMMER_MASH"
	HAMMER_INTEGER_LITERAL spansql.IntegerLiteral = 0
)

func Walk(node interface{}, f func(interface{}) interface{}) interface{} {
	switch v := node.(type) {
	case spansql.Query:
		//if err := f(v); err != nil {
		//	return err
		//}
		nn := Walk(v.Select, f)
		v.Select = nn.(spansql.Select)
		return v
	case spansql.Select:
		nn := Walk(v.Where, f)
		v.Where = nn.(spansql.BoolExpr)
		return v
	case spansql.LogicalOp:
		//if err := f(v); err != nil {
		//	return err
		//}
		{
			nn := Walk(v.LHS, f)
			if nn != nil {
				v.LHS = nn.(spansql.BoolExpr)
			}
		}
		{
			nn := Walk(v.RHS, f)
			if nn != nil {
				v.RHS = nn.(spansql.BoolExpr)
			}
		}
		return v
	case spansql.ComparisonOp:
		//if err := f(v); err != nil {
		//	return err
		//}
		{
			nn := Walk(v.LHS, f)
			if nn != nil {
				v.LHS = nn.(spansql.Expr)
			}
		}
		{
			nn := Walk(v.RHS, f)
			if nn != nil {
				v.RHS = nn.(spansql.Expr)
			}
		}
		{
			nn := Walk(v.RHS2, f)
			if nn != nil {
				v.RHS2 = nn.(spansql.Expr)
			}
		}
		return v
	case spansql.StringLiteral:
		return f(v)
	case spansql.IntegerLiteral:
		return f(v)
	default:
		return v
	}
}

//func Walk(node interface{}, f func(interface{}) interface{}) error {
//	{
//		v, ok := node.(*spansql.Query)
//		if ok {
//			if err := f(&v); err != nil {
//				return err
//			}
//			if err := Walk(&v.Select, f); err != nil {
//				return err
//			}
//			return nil
//		}
//	}
//	{
//		v, ok := node.(*spansql.Select)
//		if ok {
//			if err := f(&v); err != nil {
//				return err
//			}
//			if err := Walk(&v.Where, f); err != nil {
//				return err
//			}
//			return nil
//		}
//	}
//	{
//		v, ok := node.(*spansql.LogicalOp)
//		if ok {
//			if err := f(&v); err != nil {
//				return err
//			}
//			if err := Walk(&v.LHS, f); err != nil {
//				return err
//			}
//			if err := Walk(&v.RHS, f); err != nil {
//				return err
//			}
//		}
//	}
//	{
//		v, ok := node.(*spansql.ComparisonOp)
//		if ok {
//			if err := f(&v); err != nil {
//				return err
//			}
//			if err := Walk(&v.LHS, f); err != nil {
//				return err
//			}
//			if err := Walk(&v.RHS, f); err != nil {
//				return err
//			}
//			if err := Walk(&v.RHS2, f); err != nil {
//				return err
//			}
//		}
//	}
//
//	return nil
//}
