package analyzer

import (
	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	// Analyzer names must be valid Go identifiers, so omit the product name's hyphen.
	Name: "sealgo",
	Doc:  "check for resources that are acquired but not properly cleaned up",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	checkHTTPResponseBody(pass)
	checkFiles(pass)
	for _, file := range pass.Files {
		checkContextDeadline(pass, file)
		checkContextTimeOut(pass, file)
	}
	return nil, nil
}
