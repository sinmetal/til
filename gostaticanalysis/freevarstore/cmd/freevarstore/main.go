package main

import (
	"freevarstore"
	"golang.org/x/tools/go/analysis/unitchecker"
)

func main() { unitchecker.Main(freevarstore.Analyzer) }
