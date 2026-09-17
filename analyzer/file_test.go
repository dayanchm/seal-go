package analyzer_test

import (
	"cleanupcheck/analyzer"
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestFileClose(t *testing.T) {
	analysistest.Run(t, analysistest.TestData(), analyzer.Analyzer, "fileclose")

}
