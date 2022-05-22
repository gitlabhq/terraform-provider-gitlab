// Program tfproviderlint-plugin is a custom linter plugin for golangci-lint which runs the
// tfproviderlint analyzers.
//
// See: https://golangci-lint.run/contributing/new-linters/#create-a-plugin
package main

import (
	"github.com/bflad/tfproviderlint/passes"
	"github.com/bflad/tfproviderlint/xpasses"
	"golang.org/x/tools/go/analysis"
)

var excludes = []string{
	"XAT001",
	"XR003",
	"XS002",
}

type analyzerPlugin struct{}

func (*analyzerPlugin) GetAnalyzers() []*analysis.Analyzer {
	excludesSet := make(map[string]struct{}, len(excludes))

	for _, exclude := range excludes {
		excludesSet[exclude] = struct{}{}
	}

	var analyzers []*analysis.Analyzer

	for _, analyzer := range append(passes.AllChecks, xpasses.AllChecks...) {
		if _, isExcluded := excludesSet[analyzer.Name]; !isExcluded {
			analyzers = append(analyzers, analyzer)
		}
	}

	return analyzers
}

var AnalyzerPlugin analyzerPlugin
