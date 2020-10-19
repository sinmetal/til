package main

import (
	"github.com/gostaticanalysis/codegen/singlegenerator"
	"github.com/sinmetal/til/gostaticanalysis/genplay"
)

func main() {
	singlegenerator.Main(genplay.Generator)
}
