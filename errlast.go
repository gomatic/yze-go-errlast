// Package errlast provides a go/analysis analyzer enforcing the gomatic Go idiom
// that error is always the last return value.
package errlast

import (
	"go/ast"
	"go/types"

	goyze "github.com/gomatic/go-yze"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const message = "error must be the last return value"

// Analyzer reports error results that are not the last return value.
var Analyzer = &analysis.Analyzer{
	Name:     "errlast",
	Doc:      "reports error return values that are not last",
	Requires: []*analysis.Analyzer{inspect.Analyzer},
	Run:      run,
}

// Registration declares this analyzer to the yze framework.
var Registration = goyze.Registration{
	Name:       "errlast",
	Categories: []goyze.Category{"errors"},
	URL:        "https://docs.gomatic.dev/yze/errlast",
	Analyzer:   Analyzer,
}

// run reports each error result that is not the last return value. It inspects
// every function signature — declarations, methods, interface methods, function
// literals, and function-typed definitions — because the error-last idiom is a
// contract on any signature returning an error.
func run(pass *analysis.Pass) (any, error) {
	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	insp.Preorder([]ast.Node{(*ast.FuncType)(nil)}, func(n ast.Node) {
		checkResults(pass, n.(*ast.FuncType).Results)
	})
	return nil, nil
}

// checkResults reports an error result that is not in the final position.
func checkResults(pass *analysis.Pass, results *ast.FieldList) {
	if results == nil {
		return
	}
	positions := flattenTypes(results.List)
	last := len(positions) - 1
	for i, typ := range positions {
		if i != last && isError(pass, typ) {
			pass.Reportf(typ.Pos(), message)
		}
	}
}

// flattenTypes expands fields into one type expression per result position.
func flattenTypes(fields []*ast.Field) []ast.Expr {
	var positions []ast.Expr
	for _, field := range fields {
		count := len(field.Names)
		if count == 0 {
			count = 1
		}
		for range count {
			positions = append(positions, field.Type)
		}
	}
	return positions
}

// isError reports whether expr names the builtin error type.
func isError(pass *analysis.Pass, expr ast.Expr) bool {
	named, ok := pass.TypesInfo.TypeOf(expr).(*types.Named)
	if !ok || named.Obj().Pkg() != nil {
		return false
	}
	return named.Obj().Name() == "error"
}
