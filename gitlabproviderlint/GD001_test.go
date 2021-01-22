package main

import (
	"path/filepath"
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestGD001(t *testing.T) {
	analysistest.Run(t, analysistest.TestData(), analyzerGD001, filepath.Join("GD001", "gitlab"))
}
