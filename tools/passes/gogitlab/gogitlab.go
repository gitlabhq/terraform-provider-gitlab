package gogitlab

import (
	"fmt"
	"go/ast"
	"go/token"
	"path/filepath"
	"reflect"

	"github.com/gitlabhq/terraform-provider-gitlab/tools/passes"
	"golang.org/x/tools/go/analysis"
)

const goGitLabPackagePath = "github.com/xanzy/go-gitlab"

var Analyzer = &analysis.Analyzer{
	Doc:        "Intermediate analyzer for extracting summary data from the go-gitlab package",
	Name:       "gogitlab",
	ResultType: reflect.TypeOf((*Result)(nil)),
	// Using Facts causes the analyzer to visit dependencies;
	// otherwise it would not analyze the go-gitlab package.
	FactTypes: []analysis.Fact{
		(*typeFact)(nil),
		(*funcFact)(nil),
		(*methodFact)(nil),
		(*fieldFact)(nil),
	},
	Run: run,
}

type Result struct {
	TypeToFilenames   MultiMap
	FuncToFilenames   MultiMap
	MethodToFilenames MultiMap
	FieldToFilenames  MultiMap
}

type MultiMap map[string][]string

func run(pass *analysis.Pass) (interface{}, error) {
	if pass.Pkg.Path() == goGitLabPackagePath && !passes.IsTestPackage(pass) {
		exportTypeFacts(pass)
		exportFuncFacts(pass)
		exportMethodFacts(pass)
		exportFieldFacts(pass)
	}

	return makeResult(pass), nil
}

func exportTypeFacts(pass *analysis.Pass) {
	export := factExporter(pass, func(nameInFile nameInFile) analysis.Fact {
		fact := typeFact(nameInFile)
		return &fact
	})

	for _, file := range pass.Files {
		for _, decl := range file.Decls {
			if gen, ok := decl.(*ast.GenDecl); ok {
				switch gen.Tok {
				case token.CONST:
					for _, s := range gen.Specs {
						for _, n := range s.(*ast.ValueSpec).Names {
							export(file, n)
						}
					}
				case token.TYPE:
					for _, s := range gen.Specs {
						export(file, s.(*ast.TypeSpec).Name)
					}
				}
			}
		}
	}
}

func exportFuncFacts(pass *analysis.Pass) {
	export := factExporter(pass, func(nameInFile nameInFile) analysis.Fact {
		fact := funcFact(nameInFile)
		return &fact
	})

	for _, file := range pass.Files {
		for _, decl := range file.Decls {
			if funcDecl, ok := decl.(*ast.FuncDecl); ok {
				if funcDecl.Recv == nil {
					export(file, funcDecl.Name)
				}
			}
		}
	}
}

func exportMethodFacts(pass *analysis.Pass) {
	export := factExporter(pass, func(nameInFile nameInFile) analysis.Fact {
		fact := methodFact(nameInFile)
		return &fact
	})

	for _, file := range pass.Files {
		for _, decl := range file.Decls {
			if funcDecl, ok := decl.(*ast.FuncDecl); ok && funcDecl.Recv != nil {
				export(file, funcDecl.Name)
			}
		}
	}
}

func exportFieldFacts(pass *analysis.Pass) {
	export := factExporter(pass, func(nameInFile nameInFile) analysis.Fact {
		fact := fieldFact(nameInFile)
		return &fact
	})

	for _, file := range pass.Files {
		for _, decl := range file.Decls {
			if gen, ok := decl.(*ast.GenDecl); ok && gen.Tok == token.TYPE {
				for _, spec := range gen.Specs {
					ast.Inspect(spec, func(n ast.Node) bool {
						switch x := n.(type) {
						case *ast.TypeSpec:
							if !x.Name.IsExported() {
								return false
							}
						case *ast.Field:
							for _, name := range x.Names {
								export(file, name)
							}
						}
						return true
					})
				}
			}
		}
	}
}

func makeResult(pass *analysis.Pass) *Result {
	result := &Result{
		TypeToFilenames:   make(MultiMap),
		FuncToFilenames:   make(MultiMap),
		MethodToFilenames: make(MultiMap),
		FieldToFilenames:  make(MultiMap),
	}

	addName := func(name, filename string, m MultiMap) {
		for _, seenFilename := range m[name] {
			if filename == seenFilename {
				return
			}
		}
		m[name] = append(m[name], filename)
	}

	for _, f := range pass.AllObjectFacts() {
		switch x := f.Fact.(type) {
		case *typeFact:
			addName(x.name, x.filename, result.TypeToFilenames)
		case *funcFact:
			addName(x.name, x.filename, result.FuncToFilenames)
		case *methodFact:
			addName(x.name, x.filename, result.MethodToFilenames)
		case *fieldFact:
			addName(x.name, x.filename, result.FieldToFilenames)
		default:
			panic(fmt.Sprintf("unhandled fact %#v", f.Fact))
		}
	}

	return result
}

func factExporter(pass *analysis.Pass, newFact func(nameInFile) analysis.Fact) func(*ast.File, *ast.Ident) {
	return func(file *ast.File, ident *ast.Ident) {
		if ident.IsExported() {
			pass.ExportObjectFact(
				pass.TypesInfo.Defs[ident],
				newFact(nameInFile{
					name:     ident.Name,
					filename: filepath.Base(pass.Fset.File(file.Pos()).Name()),
				}),
			)
		}
	}
}

type nameInFile struct {
	name     string
	filename string
}

type typeFact nameInFile

func (f *typeFact) AFact() {}

type funcFact nameInFile

func (f *funcFact) AFact() {}

type methodFact nameInFile

func (f *methodFact) AFact() {}

type fieldFact nameInFile

func (f *fieldFact) AFact() {}
