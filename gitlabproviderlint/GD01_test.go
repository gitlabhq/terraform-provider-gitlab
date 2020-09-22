package main

import (
	"golang.org/x/tools/go/analysis/analysistest"
	"path/filepath"
	"testing"
)

func TestGD01(t *testing.T) {
	analysistest.Run(t, analysistest.TestData(), analyzerGD01, filepath.Join("GD01", "gitlab"))
}
