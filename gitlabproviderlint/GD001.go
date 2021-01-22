package main

import (
	"github.com/bflad/tfproviderlint/passes/commentignore"
	"golang.org/x/tools/go/analysis"
)

var analyzerGD001 = &analysis.Analyzer{
	Name: "GD001",
	Doc:  "check for resources without a docs page",
	Requires: []*analysis.Analyzer{
		analyzerParseDocs,
		analyzerProviderIndex,
		commentignore.Analyzer,
	},
	Run: analyzerGD001Run,
}

func analyzerGD001Run(pass *analysis.Pass) (interface{}, error) {
	ignorerResult := pass.ResultOf[commentignore.Analyzer].(*commentignore.Ignorer)
	docsResult := pass.ResultOf[analyzerParseDocs].(*analyzerParseDocsResult)
	providerIndexResult := pass.ResultOf[analyzerProviderIndex].(*analyzerProviderIndexResult)

	for name, node := range providerIndexResult.dataSourceNames {
		if ignorerResult.ShouldIgnore(pass.Analyzer.Name, node) {
			continue
		}
		if _, ok := docsResult.dataSources[name]; !ok {
			pass.Reportf(node.ValuePos, "Data source %q is missing a docs page", name)
		}
	}

	for name, node := range providerIndexResult.resourceNames {
		if ignorerResult.ShouldIgnore(pass.Analyzer.Name, node) {
			continue
		}
		if _, ok := docsResult.resources[name]; !ok {
			pass.Reportf(node.ValuePos, "Resource %q is missing a docs page", name)
		}
	}

	return nil, nil
}
