package analyzer_test

import (
	"testing"

	"seal-go/analyzer"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestContextDeadline(t *testing.T) {
	analysistest.Run(t, analysistest.TestData(), analyzer.Analyzer, "contextdeadline")
}
