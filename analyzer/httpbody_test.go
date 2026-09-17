package analyzer_test

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
	"seal-go/analyzer"
)

func TestHTTPResponseBody(t *testing.T) {
	analysistest.Run(t, analysistest.TestData(), analyzer.Analyzer, "httpbody")
}
