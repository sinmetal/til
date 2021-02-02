package main

import (
	"golang.org/x/tools/go/analysis/unitchecker"
	"ssasample"
)

func main() { unitchecker.Main(ssasample.Analyzer) }
