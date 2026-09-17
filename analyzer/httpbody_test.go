package analyzer_test

import (
	"testing"

	"cleanupcheck/analyzer"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestHTTPResponseBody(t *testing.T) {
	analysistest.Run(t, analysistest.TestData(), analyzer.Analyzer, "httpbody")
}
