package freevarstore

import (
	"fmt"
	"golang.org/x/tools/go/ssa"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/buildssa"
)

const doc = "freevarstore is ..."

// Analyzer is ...
var Analyzer = &analysis.Analyzer{
	Name: "freevarstore",
	Doc:  doc,
	Run:  run,
	Requires: []*analysis.Analyzer{
		buildssa.Analyzer,
	},
}

func run(pass *analysis.Pass) (interface{}, error) {
	s := pass.ResultOf[buildssa.Analyzer].(*buildssa.SSA)
	for _, f := range s.SrcFuncs {
		fmt.Println(f)
		for _, fv := range f.FreeVars {
			fmt.Printf("\tFreeVar %+v\n", fv)
		}

		for _, b := range f.Blocks {
			fmt.Printf("\tBlock %d\n", b.Index)
			for _, instr := range b.Instrs {
				fmt.Printf("\t\t%[1]T\t%[1]v(%[1]p)\n", instr)
				store, _ := instr.(*ssa.Store)
				if store == nil {
					continue // store命令以外はおしまい
				}

				freevar, _ := (store.Addr).(*ssa.FreeVar)
				if freevar != nil {
					pass.Reportf(store.Pos(), "NG")
				}
			}
		}
	}
	return nil, nil
}
