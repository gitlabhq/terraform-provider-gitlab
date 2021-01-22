package main

import (
	"go/ast"
	"reflect"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var analyzerProviderIndex = &analysis.Analyzer{
	Name:       "providerindex",
	Doc:        "parse main provider configuration for later passes",
	Requires:   []*analysis.Analyzer{inspect.Analyzer},
	Run:        analyzerProviderIndexRun,
	ResultType: reflect.TypeOf(&analyzerProviderIndexResult{}),
}

type analyzerProviderIndexResult struct {
	dataSourceNames map[string]*ast.BasicLit
	resourceNames   map[string]*ast.BasicLit
}

func analyzerProviderIndexRun(pass *analysis.Pass) (interface{}, error) {
	var result analyzerProviderIndexResult

	inspectResult := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{(*ast.KeyValueExpr)(nil)}
	inspectResult.Nodes(nodeFilter, func(n ast.Node, push bool) (prune bool) {
		if result.resourceNames != nil && result.dataSourceNames != nil {
			return false
		}

		keyValueExpr := n.(*ast.KeyValueExpr)
		if ident, ok := keyValueExpr.Key.(*ast.Ident); ok {
			switch ident.Name {
			case "ResourcesMap":
				result.resourceNames = collectMapKeys(keyValueExpr.Value)
			case "DataSourcesMap":
				result.dataSourceNames = collectMapKeys(keyValueExpr.Value)
			}
		}

		return true
	})

	return &result, nil
}

func collectMapKeys(node ast.Expr) map[string]*ast.BasicLit {
	compositeLit, ok := node.(*ast.CompositeLit)
	if !ok {
		return nil
	}

	result := make(map[string]*ast.BasicLit)
	for _, elt := range compositeLit.Elts {
		keyValueExpr, ok := elt.(*ast.KeyValueExpr)
		if !ok {
			continue
		}

		keyBasicLit, ok := keyValueExpr.Key.(*ast.BasicLit)
		if !ok {
			continue
		}

		name := strings.Trim(keyBasicLit.Value, `\"`)
		result[name] = keyBasicLit
	}

	return result
}
