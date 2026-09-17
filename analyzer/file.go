package analyzer

import (
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/analysis"
)

func isFileOpenCall(pass *analysis.Pass, call *ast.CallExpr) bool {
	selector, ok := call.Fun.(*ast.SelectorExpr)

	if !ok {
		return false
	}

	fn, ok := pass.TypesInfo.ObjectOf(selector.Sel).(*types.Func)

	if !ok || fn.Pkg() == nil || fn.Pkg().Path() != "os" {
		return false
	}

	sig, ok := fn.Type().(*types.Signature)

	if !ok || sig.Recv() != nil {
		return false
	}

	switch fn.Name() {
	case "Open", "Create", "OpenFile":
		return true
	default:
		return false
	}
}
