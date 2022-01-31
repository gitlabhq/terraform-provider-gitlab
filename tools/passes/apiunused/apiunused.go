package apiunused

import (
	"encoding/json"
	"io"
	"os"
	"reflect"

	"github.com/gitlabhq/terraform-provider-gitlab/tools/passes"
	"github.com/gitlabhq/terraform-provider-gitlab/tools/passes/gogitlab"
	"github.com/gitlabhq/terraform-provider-gitlab/tools/passes/usage"
	"golang.org/x/tools/go/analysis"
)

var Output io.Writer = io.Discard

var Analyzer = &analysis.Analyzer{
	Doc:        "Estimate the unused pieces of the go-gitlab package",
	Name:       "apiunused",
	ResultType: reflect.TypeOf((*Result)(nil)),
	Requires:   []*analysis.Analyzer{gogitlab.Analyzer, usage.Analyzer},
	Run:        run,
}

type Result struct {
	UnusedByFile map[string]Unused
}

type Unused struct {
	Types   []string `json:"types,omitempty"`
	Funcs   []string `json:"funcs,omitempty"`
	Methods []string `json:"methods,omitempty"`
	Fields  []string `json:"fields,omitempty"`
}

func run(pass *analysis.Pass) (interface{}, error) {
	goGitLab := pass.ResultOf[gogitlab.Analyzer].(*gogitlab.Result)
	usage := pass.ResultOf[usage.Analyzer].(*usage.Result)

	result := &Result{
		UnusedByFile: make(map[string]Unused),
	}

	processUnused(result, goGitLab.TypeToFilenames, usage.Types,
		func(item Unused, name string) Unused {
			item.Types = append(item.Types, name)
			return item
		})

	processUnused(result, goGitLab.FuncToFilenames, usage.Funcs,
		func(item Unused, name string) Unused {
			item.Funcs = append(item.Funcs, name)
			return item
		})

	processUnused(result, goGitLab.MethodToFilenames, usage.Methods,
		func(item Unused, name string) Unused {
			item.Methods = append(item.Methods, name)
			return item
		})

	processUnused(result, goGitLab.FieldToFilenames, usage.Fields,
		func(item Unused, name string) Unused {
			item.Fields = append(item.Fields, name)
			return item
		})

	if !passes.IsTestPackage(pass) {
		writeOutput(result)
	}

	return result, nil
}

func processUnused(result *Result, nameToFilenames gogitlab.MultiMap, seen usage.Set, mutFn func(Unused, string) Unused) {
	for name, filenames := range nameToFilenames {
		if !seen[name] {
			for _, filename := range filenames {
				result.UnusedByFile[filename] = mutFn(result.UnusedByFile[filename], name)
			}
		}
	}
}

func writeOutput(result *Result) {
	e := json.NewEncoder(os.Stdout)
	e.SetIndent("", "  ")

	_ = e.Encode(result.UnusedByFile)
}
