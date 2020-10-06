package memefish

import (
	"fmt"

	"github.com/MakeNowJust/memefish/pkg/ast"
)

func Walk(node ast.Node, f func(ast.Node) error) error {
	next := func(node ast.Node) error {
		if node == nil {
			return nil
		}
		if err := Walk(node, f); err != nil {
			return fmt.Errorf("failed walk %+v: %w", node, err)
		}
		return nil
	}

	switch v := node.(type) {
	case *ast.AddColumn:
		if err := f(v); err != nil {
			return fmt.Errorf("failed walk processing %+v: %w", v, err)
		}
		return Walk(v.Column, f)
	case *ast.Alias:
	case *ast.AlterColumn:
	case *ast.AlterColumnSet:
	case *ast.AlterTable:
	case *ast.Arg:
	case *ast.ArrayLiteral:
	case *ast.ArraySchemaType:
	case *ast.ArraySubQuery:
	case *ast.ArrayType:
	case *ast.AsAlias:
	case *ast.AtTimeZone:
	case *ast.BetweenExpr:
	case *ast.BinaryExpr:
		if err := f(v); err != nil {
			return fmt.Errorf("failed walk processing %+v: %w", v, err)
		}

		if err := next(v.Left); err != nil {
			return err
		}

		if err := next(v.Right); err != nil {
			return err
		}
	case *ast.BoolLiteral:
	case *ast.BytesLiteral:
	case *ast.CallExpr:
	case *ast.CaseElse:
	case *ast.CaseExpr:
	case *ast.CaseWhen:
	case *ast.CastExpr:
	case *ast.CastIntValue:
	case *ast.CastNumValue:
	case *ast.Cluster:
	case *ast.Collate:
	case *ast.ColumnDef:
	case *ast.ColumnDefOptions:
	case *ast.CompoundQuery:
	case *ast.CountStarExpr:
	case *ast.CreateDatabase:
	case *ast.CreateIndex:
	case *ast.CreateTable:
	case *ast.DateLiteral:
	case *ast.DefaultExpr:
	case *ast.Delete:
	case *ast.DotStar:
	case *ast.DropColumn:
	case *ast.DropIndex:
	case *ast.DropTable:
	case *ast.ExistsSubQuery:
	case *ast.ExprSelectItem:
		if err := f(v); err != nil {
			return fmt.Errorf("failed walk processing %+v: %w", v, err)
		}
		if err := next(v.Expr); err != nil {
			return err
		}
	case *ast.ExtractExpr:
	case *ast.FloatLiteral:
	case *ast.From:
	case *ast.GroupBy:
	case *ast.Having:
	case *ast.Hint:
	case *ast.HintRecord:
	case *ast.Ident:
	case *ast.InExpr:
	case *ast.IndexExpr:
	case *ast.IndexKey:
	case *ast.Insert:
	case *ast.IntLiteral:
		if err := f(v); err != nil {
			return fmt.Errorf("failed walk processing %+v: %w", v, err)
		}
	case *ast.InterleaveIn:
	case *ast.IsBoolExpr:
	case *ast.IsNullExpr:
	case *ast.Join:
	case *ast.Limit:
	case *ast.NullLiteral:
	case *ast.Offset:
	case *ast.On:
	case *ast.OrderBy:
	case *ast.OrderByItem:
	case *ast.Param:
	case *ast.ParenExpr:
	case *ast.ParenTableExpr:
	case *ast.Path:
	case *ast.QueryStatement:
		if err := f(v); err != nil {
			return err
		}
		if err := next(v.Query); err != nil {
			return err
		}
	case *ast.ScalarSchemaType:
	case *ast.ScalarSubQuery:
	case *ast.Select:
		if err := f(v); err != nil {
			return fmt.Errorf("failed walk processing %+v: %w", v, err)
		}

		for _, result := range v.Results {
			if err := f(result); err != nil {
				return fmt.Errorf("failed walk processing %+v: %w", result, err)
			}
		}

		if err := next(v.From); err != nil {
			return err
		}
		if err := next(v.Where); err != nil {
			return err
		}
		if err := next(v.GroupBy); err != nil {
			return err
		}
		if err := next(v.Having); err != nil {
			return err
		}
		if err := next(v.OrderBy); err != nil {
			return err
		}
	case *ast.SelectorExpr:
	case *ast.SetOnDelete:
	case *ast.SimpleType:
	case *ast.SizedSchemaType:
	case *ast.Star:
	case *ast.Storing:
	case *ast.StringLiteral:
		if err := f(v); err != nil {
			return fmt.Errorf("failed walk processing %+v: %w", v, err)
		}
	case *ast.StructField:
	case *ast.StructLiteral:
	case *ast.StructType:
	case *ast.SubQuery:
	case *ast.SubQueryInCondition:
	case *ast.SubQueryInput:
	case *ast.SubQueryTableExpr:
	case *ast.TableName:
	case *ast.TableSample:
	case *ast.TableSampleSize:
	case *ast.TimestampLiteral:
	case *ast.UnaryExpr:
	case *ast.Unnest:
	case *ast.UnnestInCondition:
	case *ast.Update:
	case *ast.UpdateItem:
	case *ast.Using:
	case *ast.ValuesInCondition:
	case *ast.ValuesInput:
	case *ast.ValuesRow:
	case *ast.Where:
		if err := f(v); err != nil {
			return fmt.Errorf("failed walk processing %+v: %w", v, err)
		}

		if err := next(v.Expr); err != nil {
			return err
		}
	case *ast.WithOffset:
	default:
		return fmt.Errorf("unsupported type %T %+v: ", node, node)
	}
	return nil
}
