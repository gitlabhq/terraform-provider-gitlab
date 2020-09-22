package main

import (
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
)

var allAnalyzers = []*analysis.Analyzer{
	analyzerGD001,
}

func main() {
	multichecker.Main(allAnalyzers...)
}
