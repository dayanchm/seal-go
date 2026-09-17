package analyzer

import (
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/analysis"
)

func checkContextCancellation(pass *analysis.Pass, file *ast.File) {
	ast.Inspect(file, func(n ast.Node) bool {

		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}
		sel, ok := call.Fun.(*ast.SelectorExpr)
		if !ok {
			return true
		}
		pkg, ok := sel.X.(*ast.Ident)
		if !ok {
			return true
		}
		if pkg.Name == "context" && sel.Sel.Name == "WithCancel" {
			pass.Reportf(call.Pos(), "context.WithCancel requires cancel()")
		}
		return true
	})
}

func checkContextTimeOut(pass *analysis.Pass, file *ast.File) {
	ast.Inspect(file, func(n ast.Node) bool {
		assign, ok := n.(*ast.AssignStmt)
		if !ok {
			return true
		}

		if len(assign.Rhs) != 1 || len(assign.Lhs) != 2 {
			return true
		}
		call, ok := assign.Rhs[0].(*ast.CallExpr)
		if !ok {
			return true
		}

		sel, ok := call.Fun.(*ast.SelectorExpr)
		if !ok {
			return true
		}

		pkg, ok := sel.X.(*ast.Ident)
		if !ok {
			return true
		}
		if pkg.Name != "context" || sel.Sel.Name != "WithTimeout" {
			return true
		}
		cancelIdent, ok := assign.Lhs[1].(*ast.Ident)
		if !ok {
			return true
		}

		if cancelIdent.Name == "_" {
			pass.Reportf(
				cancelIdent.Pos(),
				"context.WithTimeout cancellation function is discarded",
			)
		}

		return true
	})

}

func checkContextDeadline(pass *analysis.Pass, file *ast.File) {
	ast.Inspect(file, func(n ast.Node) bool {
		var cancelIdent *ast.Ident
		var rhs ast.Expr

		switch node := n.(type) {
		case *ast.AssignStmt:
			if len(node.Rhs) != 1 || len(node.Lhs) != 2 {
				return true
			}
			cancelIdent, _ = node.Lhs[1].(*ast.Ident)
			rhs = node.Rhs[0]
		case *ast.ValueSpec:
			if len(node.Values) != 1 || len(node.Names) != 2 {
				return true
			}
			cancelIdent = node.Names[1]
			rhs = node.Values[0]
		default:
			return true
		}

		if cancelIdent == nil || cancelIdent.Name != "_" {
			return true
		}

		call, ok := ast.Unparen(rhs).(*ast.CallExpr)
		if !ok {
			return true
		}

		var callee *ast.Ident
		switch fun := ast.Unparen(call.Fun).(type) {
		case *ast.SelectorExpr:
			callee = fun.Sel
		case *ast.Ident:
			callee = fun
		default:
			return true
		}

		fn, ok := pass.TypesInfo.ObjectOf(callee).(*types.Func)
		if !ok || fn.Pkg() == nil || fn.Pkg().Path() != "context" || fn.Name() != "WithDeadline" {
			return true
		}

		pass.Reportf(cancelIdent.Pos(), "context.WithDeadline cancellation function is discarded")
		return true
	})
}
