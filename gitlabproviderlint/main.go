package main

import (
	tfpasses "github.com/bflad/tfproviderlint/passes"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
)

var allAnalyzers = []*analysis.Analyzer{
	analyzerGD001,
}

func main() {
	// Add the standard tfproviderlint analyzers.
	allAnalyzers = append(allAnalyzers, tfpasses.AllChecks...)

	// Run all analyzers.
	multichecker.Main(allAnalyzers...)
}
