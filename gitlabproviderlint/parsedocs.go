package main

import (
	"github.com/yuin/goldmark"
	goldmarkast "github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
	"golang.org/x/tools/go/analysis"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
)

var analyzerParseDocs = &analysis.Analyzer{
	Name:       "parsedocs",
	Doc:        "parse provider documentation for later passes",
	Run:        analyzerParseDocsRun,
	ResultType: reflect.TypeOf(&analyzerParseDocsResult{}),
}

type analyzerParseDocsResult struct {
	dataSources []analyzerParseDocsResultPage
	resources   []analyzerParseDocsResultPage
}

type analyzerParseDocsResultPage struct {
	filename string
	content  goldmarkast.Node
}

func analyzerParseDocsRun(pass *analysis.Pass) (interface{}, error) {
	var result analyzerParseDocsResult

	docsPath := discoverDocsPath(pass)
	if docsPath == "" {
		return &result, nil
	}

	err := filepath.Walk(docsPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if filepath.Ext(path) != ".md" {
			return nil
		}

		switch filepath.Dir(path) {
		case filepath.Join(docsPath, "resources"):
			page, err := newAnalyzerParseDocsPage(path)
			if err != nil {
				return err
			}
			result.resources = append(result.resources, page)
		case filepath.Join(docsPath, "data-sources"):
			page, err := newAnalyzerParseDocsPage(path)
			if err != nil {
				return err
			}
			result.dataSources = append(result.dataSources, page)
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

func newAnalyzerParseDocsPage(path string) (analyzerParseDocsResultPage, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return analyzerParseDocsResultPage{}, err
	}

	content := goldmark.DefaultParser().Parse(text.NewReader(b))

	return analyzerParseDocsResultPage{
		filename: filepath.Base(path),
		content:  content,
	}, nil
}
