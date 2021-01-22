package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/yuin/goldmark"
	goldmarkast "github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
	"golang.org/x/tools/go/analysis"
)

var analyzerParseDocs = &analysis.Analyzer{
	Name:       "parsedocs",
	Doc:        "parse provider documentation for later passes",
	Run:        analyzerParseDocsRun,
	ResultType: reflect.TypeOf(&analyzerParseDocsResult{}),
}

type analyzerParseDocsResult struct {
	dataSources map[string]analyzerParseDocsResultPage
	resources   map[string]analyzerParseDocsResultPage
}

type analyzerParseDocsResultPage struct {
	filepath string
	source   []byte
	tree     goldmarkast.Node
}

func analyzerParseDocsRun(pass *analysis.Pass) (interface{}, error) {
	result := analyzerParseDocsResult{
		dataSources: make(map[string]analyzerParseDocsResultPage),
		resources:   make(map[string]analyzerParseDocsResultPage),
	}

	docsPath := discoverDocsPath(pass)
	if docsPath == "" {
		return &result, nil
	}

	projectBasePath := filepath.Dir(docsPath)

	err := filepath.Walk(docsPath, func(pageAbsPath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if filepath.Ext(pageAbsPath) != ".md" {
			return nil
		}

		name := "gitlab_" + strings.TrimSuffix(filepath.Base(pageAbsPath), ".md")

		switch filepath.Dir(pageAbsPath) {
		case filepath.Join(docsPath, "resources"):
			page, err := newAnalyzerParseDocsPage(pageAbsPath, projectBasePath)
			if err != nil {
				return err
			}
			result.resources[name] = page
		case filepath.Join(docsPath, "data-sources"):
			page, err := newAnalyzerParseDocsPage(pageAbsPath, projectBasePath)
			if err != nil {
				return err
			}
			result.dataSources[name] = page
		}

		return nil
	})

	return &result, err
}

// discoverDocsPath returns the absolute filepath for the docs.
// The reason we cannot simply hardcode the path is because this analyzer also needs to run in tests.
func discoverDocsPath(pass *analysis.Pass) string {
	// Get the absolute filepath for the package being analyzed.

	if len(pass.Files) == 0 {
		return ""
	}

	fileToken := pass.Fset.File(pass.Files[0].Package)
	if fileToken == nil {
		return ""
	}

	// Walk backwards from the package to find the docs path.

	lastPackage := ""
	thisPackage := filepath.Dir(fileToken.Name())
	maxDepth := 10

	// Loop until there is no parent dir. The maxDepth is a safeguard.
	for depth := 0; lastPackage != thisPackage && depth < maxDepth; depth++ {
		files, err := ioutil.ReadDir(thisPackage)
		if err != nil {
			return ""
		}

		for _, f := range files {
			if f.IsDir() && f.Name() == "docs" {
				return filepath.Join(thisPackage, f.Name())
			}
		}

		lastPackage = thisPackage
		thisPackage = filepath.Dir(thisPackage)
	}

	return ""
}

func newAnalyzerParseDocsPage(pageAbsPath, projectBasePath string) (analyzerParseDocsResultPage, error) {
	source, err := ioutil.ReadFile(pageAbsPath)
	if err != nil {
		return analyzerParseDocsResultPage{}, err
	}

	tree := goldmark.DefaultParser().Parse(text.NewReader(source))

	relPath, err := filepath.Rel(projectBasePath, pageAbsPath)
	if err != nil {
		return analyzerParseDocsResultPage{}, err
	}

	return analyzerParseDocsResultPage{
		filepath: relPath,
		source:   source,
		tree:     tree,
	}, nil
}
