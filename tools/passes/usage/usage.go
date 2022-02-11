package usage

import (
	"go/ast"
	"reflect"

	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Doc:        "Estimate usage of the go-gitlab package",
	Name:       "apiCoverage",
	ResultType: reflect.TypeOf((*Result)(nil)),
	Run:        run,
}

type Result struct {
	Types   Set
	Funcs   Set
	Methods Set
	Fields  Set
}

type Set map[string]bool

func run(pass *analysis.Pass) (interface{}, error) {
	result := &Result{
		Types:   make(Set),
		Funcs:   make(Set),
		Methods: make(Set),
		Fields:  make(Set),
	}

	for _, f := range pass.Files {
		ast.Inspect(f, func(n ast.Node) bool {
			switch x := n.(type) {
			case *ast.SelectorExpr:
				if ident, ok := x.X.(*ast.Ident); ok {
					if ident.Name == "gitlab" {
						result.Types[x.Sel.Name] = true
						result.Funcs[x.Sel.Name] = true
						break
					}
				}
				result.Methods[x.Sel.Name] = true
				result.Fields[x.Sel.Name] = true
			case *ast.KeyValueExpr:
				if ident, ok := x.Key.(*ast.Ident); ok {
					result.Fields[ident.Name] = true
				}
			}
			return true
		})
	}

	return result, nil
}
