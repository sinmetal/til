package closeclose

import (
	"fmt"
	"golang.org/x/tools/go/ssa"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/buildssa"
)

const doc = "closeclose is ..."

// Analyzer is ...
var Analyzer = &analysis.Analyzer{
	Name: "closeclose",
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
		for _, b := range f.Blocks {
			fmt.Printf("\tBlock %d\n", b.Index)
			for _, instr := range b.Instrs {
				fmt.Printf("\t\t%[1]T\t%[1]v(%[1]p)\n", instr)
				switch instr := instr.(type) {
				case ssa.CallInstruction:
					cl, _ := instr.Common().Value.(*ssa.Builtin)
					if cl == nil || cl.Name() != "close" {
						continue
					}
					args := instr.Common().Args
					if len(args) != 1 {
						continue
					}
				}
				for _, v := range instr.Operands(nil) {
					if v != nil {
						fmt.Printf("\t\t\t%[1]T\t%[1]v(%[1]p)\n", *v)
					}
				}
			}
		}
	}
	return nil, nil
}
