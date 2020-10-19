package genplay

import (
	"bytes"
	"fmt"
	"go/format"
	"go/types"
	"os"

	"github.com/gostaticanalysis/analysisutil"
	"github.com/gostaticanalysis/codegen"
	"github.com/gostaticanalysis/knife"
)

const doc = "genplay is ..."

var (
	flagOutput string
)

func init() {
	Generator.Flags.StringVar(&flagOutput, "o", "", "output file name")
}

var Generator = &codegen.Generator{
	Name: "genplay",
	Doc:  doc,
	Run:  run,
}

func run(pass *codegen.Pass) error {
	sts := map[string]*knife.Struct{}

	s := pass.Pkg.Scope()
	for _, name := range s.Names() {
		obj := s.Lookup(name)
		if !obj.Exported() {
			continue
		}
		st, _ := analysisutil.Under(obj.Type()).(*types.Struct)
		if st != nil {
			sts[name] = knife.NewStruct(st)
		}
	}

	td := &knife.TempalteData{
		Fset:      pass.Fset,
		Files:     pass.Files,
		TypesInfo: pass.TypesInfo,
		Pkg:       pass.Pkg,
	}
	t, err := knife.NewTemplate(td).Parse(tmpl)
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, sts); err != nil {
		return err
	}

	src, err := format.Source(buf.Bytes())
	if err != nil {
		return err
	}

	if flagOutput == "" {
		pass.Print(string(src))
		return nil
	}

	f, err := os.Create(flagOutput)
	if err != nil {
		return err
	}

	fmt.Fprint(f, string(src))

	if err := f.Close(); err != nil {
		return err
	}

	return nil
}

var tmpl = `// Code generated by genplay; DO NOT EDIT.
package {{(pkg).Name}}
{{range $tn, $t := .}}
	var _ = {{$tn}}
	{{range $fn, $f := $t.Fields}}
		{{if (and (not $f.Exported) (eq $f.Tag "gen:\"getset\""))}}
			func (e *{{$tn}}) Get{{$fn}}() {{$f.Type}} {
				return this.{{$fn}}
			}
		{{end}}
	{{end}}
{{end}}
`
