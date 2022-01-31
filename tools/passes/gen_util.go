package passes

import (
	"strings"

	"golang.org/x/tools/go/analysis"
)

func IsTestPackage(pass *analysis.Pass) bool {
	if strings.HasSuffix(pass.Pkg.Path(), ".test") {
		return true
	}

	for _, f := range pass.Files {
		if strings.HasSuffix(pass.Fset.File(f.Pos()).Name(), "_test.go") {
			return true
		}
	}

	return false
}
