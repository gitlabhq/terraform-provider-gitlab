package main

import (
	"go/ast"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
	"strings"
)

var analyzerGD01 = &analysis.Analyzer{
	Name:     "GD01",
	Doc:      "check for resources without a docs page",
	Run:      analyzerGD01Run,
	Requires: []*analysis.Analyzer{analyzerParseDocs, inspect.Analyzer},
}

func analyzerGD01Run(pass *analysis.Pass) (interface{}, error) {
	docsResult := pass.ResultOf[analyzerParseDocs].(*analyzerParseDocsResult)
	resourceNames, dataSourceNames := listResourceNames(pass)

	for _, name := range resourceNames {
		if wantFilename, exists := hasDocsPage(name, docsResult.resources); !exists {
			pass.Reportf(name.ValuePos, "Resource %q is missing a docs page named %q", strings.Trim(name.Value, `\"`), wantFilename)
		}
	}

	for _, name := range dataSourceNames {
		if wantFilename, exists := hasDocsPage(name, docsResult.dataSources); !exists {
			pass.Reportf(name.ValuePos, "Data source %q is missing a docs page named %q", strings.Trim(name.Value, `\"`), wantFilename)
		}
	}

	return nil, nil
}

func listResourceNames(pass *analysis.Pass) (resourceNames, dataSourceNames []*ast.BasicLit) {
	inspectResult := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{(*ast.KeyValueExpr)(nil)}
	inspectResult.Nodes(nodeFilter, func(n ast.Node, push bool) (prune bool) {
		if resourceNames != nil && dataSourceNames != nil {
			return false
		}

		keyValueExpr := n.(*ast.KeyValueExpr)
		if ident, ok := keyValueExpr.Key.(*ast.Ident); ok {
			switch ident.Name {
			case "ResourcesMap":
				resourceNames = collectMapKeys(keyValueExpr.Value)
			case "DataSourcesMap":
				dataSourceNames = collectMapKeys(keyValueExpr.Value)
			}
		}

		return true
	})

	return resourceNames, dataSourceNames
}

func collectMapKeys(node ast.Expr) []*ast.BasicLit {
	compositeLit, ok := node.(*ast.CompositeLit)
	if !ok {
		return nil
	}

	var result []*ast.BasicLit
	for _, elt := range compositeLit.Elts {
		keyValueExpr, ok := elt.(*ast.KeyValueExpr)
		if !ok {
			continue
		}

		keyBasicLit, ok := keyValueExpr.Key.(*ast.BasicLit)
		if !ok {
			continue
		}

		result = append(result, keyBasicLit)
	}

	return result
}

func hasDocsPage(name *ast.BasicLit, pages []analyzerParseDocsResultPage) (wantFilename string, exists bool) {
	wantFilename = strings.TrimPrefix(strings.Trim(name.Value, `\"`), "gitlab_") + ".md"

	for _, page := range pages {
		if page.filename == wantFilename {
			return wantFilename, true
		}
	}

	return wantFilename, false
}
