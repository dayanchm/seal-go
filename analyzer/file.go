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

func findOpenedFiles(
	pass *analysis.Pass,
	body *ast.BlockStmt,
) map[types.Object]ast.Expr {
	files := make(map[types.Object]ast.Expr)

	ast.Inspect(body, func(node ast.Node) bool {
		if _, ok := node.(*ast.FuncLit); ok {
			return false
		}

		assign, ok := node.(*ast.AssignStmt)
		if !ok || len(assign.Rhs) != 1 || len(assign.Lhs) == 0 {
			return true
		}

		call, ok := assign.Rhs[0].(*ast.CallExpr)
		if !ok || !isFileOpenCall(pass, call) {
			return true
		}

		ident, ok := assign.Lhs[0].(*ast.Ident)
		if !ok || ident.Name == "_" {
			return true
		}

		obj := pass.TypesInfo.ObjectOf(ident)
		if obj != nil {
			files[obj] = call
		}

		return true
	})

	return files
}

func checkFunctionForFiles(
	pass *analysis.Pass,
	body *ast.BlockStmt,
) {
	files := findOpenedFiles(pass, body)
	if len(files) == 0 {
		return
	}

	closed := make(map[types.Object]bool)

	ast.Inspect(body, func(node ast.Node) bool {
		if _, ok := node.(*ast.FuncLit); ok {
			return false
		}

		call, ok := node.(*ast.CallExpr)
		if !ok {
			return true
		}

		selector, ok := call.Fun.(*ast.SelectorExpr)
		if !ok || selector.Sel.Name != "Close" {
			return true
		}

		ident, ok := selector.X.(*ast.Ident)
		if !ok {
			return true
		}

		obj := pass.TypesInfo.ObjectOf(ident)
		if _, exists := files[obj]; exists {
			closed[obj] = true
		}

		return true
	})

	for obj, expr := range files {
		if closed[obj] {
			continue
		}

		pass.Reportf(
			expr.Pos(),
			"file %q is not closed; call defer %s.Close()",
			obj.Name(),
			obj.Name(),
		)
	}
}

func checkFiles(pass *analysis.Pass) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(node ast.Node) bool {
			switch fn := node.(type) {
			case *ast.FuncDecl:
				if fn.Body != nil {
					checkFunctionForFiles(pass, fn.Body)
				}

			case *ast.FuncLit:
				checkFunctionForFiles(pass, fn.Body)
			}

			return true
		})
	}
}
