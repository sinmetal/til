package ssasample

import (
	"fmt"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/buildssa"
	"golang.org/x/tools/go/ssa"
)

const doc = "ssasample is ..."

// Analyzer is ...
var Analyzer = &analysis.Analyzer{
	Name: "ssasample",
	Doc:  doc,
	Run:  run,
	Requires: []*analysis.Analyzer{
		buildssa.Analyzer,
	},
}

// これだとクロージャの中は追いかけてない
func run(pass *analysis.Pass) (interface{}, error) {
	s := pass.ResultOf[buildssa.Analyzer].(*buildssa.SSA)
	for _, f := range s.SrcFuncs {
		fmt.Println(f)
		for _, param := range f.Params {
			fmt.Println(param, param.Referrers()) // parameter n : int &[] になってて、使ってないっぽいぞって雰囲気
		}
		fmt.Println(f.Params)
		used := make(map[ssa.Value]bool, len(f.Params))
		for _, b := range f.Blocks {
			fmt.Printf("\tBlock %d\n", b.Index)
			for _, instr := range b.Instrs {
				fmt.Printf("\t\t%[1]T\t%[1]v(%[1]p)\n", instr)
				for _, v := range instr.Operands(nil) {
					if v != nil {
						fmt.Printf("\t\t\t%[1]T\t%[1]v(%[1]p)\n", *v)
						used[*v] = true
					}
				}
			}
		}

		for _, p := range f.Params {
			if !used[p] {
				pass.Reportf(p.Pos(), "unused")
			}
		}
	}
	return nil, nil
}
