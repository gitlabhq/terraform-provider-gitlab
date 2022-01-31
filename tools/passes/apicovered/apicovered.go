package apicovered

import (
	"fmt"
	"io"
	"os"
	"reflect"
	"sort"
	"text/tabwriter"

	"github.com/gitlabhq/terraform-provider-gitlab/tools/passes"
	"github.com/gitlabhq/terraform-provider-gitlab/tools/passes/gogitlab"
	"github.com/gitlabhq/terraform-provider-gitlab/tools/passes/usage"
	"golang.org/x/tools/go/analysis"
)

var Output io.Writer = io.Discard

var Analyzer = &analysis.Analyzer{
	Doc:        "Estimate usage of the go-gitlab package",
	Name:       "apicovered",
	ResultType: reflect.TypeOf((*Result)(nil)),
	Requires:   []*analysis.Analyzer{gogitlab.Analyzer, usage.Analyzer},
	Run:        run,
}

type Result struct {
	CoverageByFile map[string]Fraction
}

type Fraction struct {
	Count int
	Total int
}

func run(pass *analysis.Pass) (interface{}, error) {
	goGitLab := pass.ResultOf[gogitlab.Analyzer].(*gogitlab.Result)
	usage := pass.ResultOf[usage.Analyzer].(*usage.Result)

	result := &Result{
		CoverageByFile: make(map[string]Fraction),
	}

	process(result, goGitLab.TypeToFilenames, usage.Types)
	process(result, goGitLab.FuncToFilenames, usage.Funcs)
	process(result, goGitLab.MethodToFilenames, usage.Methods)
	process(result, goGitLab.FieldToFilenames, usage.Fields)

	if !passes.IsTestPackage(pass) {
		writeOutput(result)
	}

	return result, nil
}

func process(result *Result, nameToFilenames gogitlab.MultiMap, seen usage.Set) {
	for _, filenames := range nameToFilenames {
		for _, filename := range filenames {
			coverage := result.CoverageByFile[filename]
			coverage.Total++
			result.CoverageByFile[filename] = coverage
		}
	}

	for name := range seen {
		for _, filename := range nameToFilenames[name] {
			coverage := result.CoverageByFile[filename]
			coverage.Count++
			result.CoverageByFile[filename] = coverage
		}
	}
}

func writeOutput(result *Result) {
	type row struct {
		filename        string
		coverageString  string
		coveragePercent int
	}

	newRow := func(filename string, fraction Fraction) row {
		return row{
			filename:        filename,
			coverageString:  fmt.Sprintf("%d/%d", fraction.Count, fraction.Total),
			coveragePercent: makePercent(fraction.Count, fraction.Total),
		}
	}

	var rows []row

	for filename, coverage := range result.CoverageByFile {
		rows = append(rows, newRow(filename, coverage))
	}

	sort.Slice(rows, func(i, j int) bool {
		if rows[i].coveragePercent == rows[j].coveragePercent {
			return rows[i].filename < rows[j].filename
		}
		return rows[i].coveragePercent < rows[j].coveragePercent
	})

	totalRow := func() row {
		var totalCoverage Fraction
		for _, coverage := range result.CoverageByFile {
			totalCoverage.Count += coverage.Count
			totalCoverage.Total += coverage.Total
		}
		return newRow("Total", totalCoverage)
	}()

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)

	writeRow := func(row row) {
		fmt.Fprintf(w, "%s\t%d%%\t%s\n", row.filename, row.coveragePercent, row.coverageString)
	}

	fmt.Fprintf(w, "Filename\tCoverage\tLines\n")
	fmt.Fprintln(w, "--------\t--------\t-----")

	for _, row := range rows {
		writeRow(row)
	}

	fmt.Fprintln(w, "---\t---\t---")
	writeRow(totalRow)

	w.Flush()
}

func makePercent(n, d int) int {
	if d == 0 {
		if n == 0 {
			return 100
		}
		return 0
	}
	return 100 * n / d
}
