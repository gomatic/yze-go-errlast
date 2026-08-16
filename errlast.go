// Package errlast provides a go/analysis analyzer enforcing the gomatic Go idiom
// that error is always the last return value.
//
// # The rule, and the two ways a result list breaks it
//
// A result list can be wrong in exactly one of two ways, so the analyzer emits at
// most ONE diagnostic per signature and each one names a remedy that is available:
//
//   - Exactly one result is an error and it is not the last one. Reported at that
//     result with [message]. The remedy is to move it, and it always exists.
//   - More than one result is an error. There is NO compliant order — at most one
//     of them can be last, whatever the author writes — so a reorder is an
//     instruction nobody can carry out. Reported once, at the first of them, with
//     [multiErrorMessage], which prescribes the two remedies that do exist —
//     joining them, or returning one value that carries them.
//
// An earlier revision reported every non-last error position with [message], so
// (error, error) drew an instruction no permutation satisfies and testdata pinned
// it as expected output. A diagnostic whose remedy cannot be taken leaves the
// author one move — silence the rule — which is the population this suite exists
// to remove, not to manufacture.
//
// # The count follows the violation, never the spelling
//
// (err, err2 error, n int) and (err error, err2 error, n int) are the same
// two-error violation and draw the same single diagnostic. A field declaring
// several names declares several results sharing ONE type expression, so an
// earlier revision that reported at that shared expression emitted N diagnostics
// at the IDENTICAL position — invisible in the text driver, which deduplicates by
// position, invisible under the runner, which collapses them while absorbing an
// unrelated package-variant duplication, and reported in full by this command's
// own -json output. Under a quota model the finding count is the charged unit, so
// a count that moves with comma placement charges unstably. Each result is
// therefore reported at its OWN identifier when it is named, and at the type
// expression only when the field is unnamed.
//
// # What counts as an error result
//
// The property is the one the rule is about: the caller receives an error value
// it has to check. So a result is in scope when, after unaliasing, its type is an
// interface type that implements the builtin error interface. The builtin error,
// an alias of it (type E = error), a defined type over it (type E error) and an
// interface embedding it (type Wrapped interface{ error }) are one property
// spelled four ways, and all four are reported.
//
// Keying on the SPELLING instead is a one-line, behaviour-preserving disablement,
// and it was one. An earlier revision accepted only a *types.Named whose object is
// universe-scope and literally named "error", so `type E error` silenced the
// identical violating shape while E and error stayed mutually assignable and
// errors.Is kept working — an escape that appears in no configuration file, costs
// the author none of the property, and takes any codebase that gives its error
// interface a domain name silently out of the rule's reach. Naming a type is not a
// design change and must not be an exemption.
//
// A package-scope `type error interface{ Error() string }`, which shadows the
// predeclared identifier, is in scope for the same reason: it is an interface and
// a value of it is an error. The revision that matched the name was silent on it.
//
// Two shapes are out of scope, and each has its own reason rather than a shared one:
//
//   - A CONCRETE type implementing error (func f() (*MyErr, int)). Its underlying
//     type is not an interface, so the property above does not hold; returning a
//     concrete error type is its own smell and not this analyzer's rule.
//   - A TYPE PARAMETER (func f[E error]() (E, int)). A type parameter's underlying
//     type is its constraint, so testing the property would put every
//     error-constrained parameter in scope by accident. What the caller receives is
//     fixed at the instantiation, not here.
//
// # Where it looks, and what that costs
//
// Every function type in the source: function and method declarations, interface
// methods, function literals, named function types, and equally a function type
// written as a parameter, a struct field, a conversion or a composite element. The
// idiom is a contract on any signature returning an error and each of those is a
// signature the author writes in their own file. An earlier revision's comment
// enumerated five sites while the code walked all of them.
//
// A signature whose result order is fixed by somebody else's contract is REPORTED,
// and there is no exemption for it. Two things are true and both are stated here
// rather than assumed:
//
//   - The usual shape has a remedy, and it is a better design than the one that
//     draws the finding. A module that mirrors a dependency's result order into its
//     own interface inherits that order and cannot reorder it without ceasing to
//     match — golang.org/x/sync/singleflight publishes (v any, err error, shared
//     bool), and reordering a local mirror of it stops the module compiling. The
//     remedy is to ADAPT at the boundary instead of mirroring: declare the module's
//     own error-last contract and let one adapter method absorb the foreign order.
//     Verified rather than argued — the adapted module compiles and this analyzer is
//     silent on it.
//   - The irreducible shape is a type implementing a FOREIGN INTERFACE whose error
//     result is not last, where no adapter helps because the adapter's own method
//     carries the foreign signature. That shape is reported and cannot be fixed.
//     There is deliberately no exemption for it, for two measured reasons. Its
//     population is ~0: scanning 7,665 files of the Go standard library found ZERO
//     such interfaces, and 204,559 files of a populated module cache found eight, of
//     which three are this analyzer's own fixtures, two are an unexported interface
//     in a 2019 x/tools snapshot, one is an unexported package-level var, and the
//     one externally importable case lives in a package its own author named
//     `private`. And every matcher that could recognise it — a name-and-signature
//     match against an imported type, a satisfied interface, a module-path test — is
//     forged in a few lines that acquire none of the property, which is the
//     shape-keyed exemption this suite has already paid twice to learn about. An
//     exemption serving a population of one, forgeable by anyone, is a disablement
//     channel with a rule attached.
package errlast

import (
	"go/ast"
	"go/token"
	"go/types"

	goyze "github.com/gomatic/go-yze"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const (
	// message is reported for the single-error violation, whose remedy is a move.
	message = "error must be the last return value"
	// multiErrorMessage is reported for a result list carrying more than one
	// error, which no ordering satisfies, so it prescribes the remedies that do
	// exist. Both are named because joining is lossy where the results are told
	// apart by POSITION rather than by sentinel — crypto/tls's `(clientErr,
	// serverErr error)` is one such shape — and a remedy that silently
	// discards what the signature was carrying is not one the author can take.
	multiErrorMessage = "a signature returning more than one error has no compliant order; " +
		"return one joined error, or one value carrying them"
)

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

// resultIndex is a result's position in the flattened result list.
type resultIndex int

// errorResult is one result whose type is an error interface: where a diagnostic
// about it is reported, and where it sits among the results.
type errorResult struct {
	pos   token.Pos
	index resultIndex
}

// run inspects every function type in the source, because the error-last idiom is
// a contract on any signature returning an error.
func run(pass *analysis.Pass) (any, error) {
	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	insp.Preorder([]ast.Node{(*ast.FuncType)(nil)}, func(n ast.Node) {
		checkResults(pass, n.(*ast.FuncType).Results)
	})
	return nil, nil
}

// checkResults reports the one way this result list breaks the rule, if it does.
//
// The nil guard is load-bearing rather than defensive: a function type with no
// results carries a nil *ast.FieldList, so reading results.List without it
// dereferences nil and the binary dies mid-analysis. Deleting it was measured, and
// it panics on every probe module; the `noResults` fixture is what enters it.
func checkResults(pass *analysis.Pass, results *ast.FieldList) {
	if results == nil {
		return
	}
	count, found := scanResults(pass, results.List)
	if len(found) > 1 {
		pass.Reportf(found[0].pos, multiErrorMessage)
		return
	}
	if len(found) == 1 && found[0].index != count-1 {
		pass.Reportf(found[0].pos, message)
	}
}

// scanResults walks the declared results in order, returning how many there are
// and which of them are error interfaces.
func scanResults(pass *analysis.Pass, fields []*ast.Field) (resultIndex, []errorResult) {
	var count resultIndex
	var found []errorResult
	for _, field := range fields {
		for _, pos := range resultPositions(field) {
			found = appendError(found, pass, field.Type, errorResult{pos: pos, index: count})
			count++
		}
	}
	return count, found
}

// resultPositions is the position of each result a field declares: the declared
// identifier for each name, or the type expression when the field is unnamed. This
// is what keeps a diagnostic's position independent of how the results are spelled.
func resultPositions(field *ast.Field) []token.Pos {
	if len(field.Names) == 0 {
		return []token.Pos{field.Type.Pos()}
	}
	positions := make([]token.Pos, 0, len(field.Names))
	for _, name := range field.Names {
		positions = append(positions, name.Pos())
	}
	return positions
}

// appendError keeps the result when its type expression is an error interface.
func appendError(
	found []errorResult,
	pass *analysis.Pass,
	typ ast.Expr,
	result errorResult,
) []errorResult {
	if !isErrorResult(pass, typ) {
		return found
	}
	return append(found, result)
}

// isErrorResult reports whether the result's type is an error interface.
//
// A type parameter is refused before the property is tested: its underlying type
// is its constraint, so an error-constrained parameter would satisfy the property
// without the caller ever receiving an interface. The nil arm shares that refusal
// because a result expression with no recorded type is not an error result either.
func isErrorResult(pass *analysis.Pass, expr ast.Expr) bool {
	switch typ := types.Unalias(pass.TypesInfo.TypeOf(expr)).(type) {
	case nil, *types.TypeParam:
		return false
	default:
		return isErrorInterface(typ)
	}
}

// isErrorInterface reports whether the type is an interface that implements the
// builtin error interface. Both halves are required: the implements test alone
// admits every concrete error type, which is out of scope with its own reason.
func isErrorInterface(typ types.Type) bool {
	_, ok := typ.Underlying().(*types.Interface)
	return ok && types.Implements(typ, errorInterface())
}

// errorInterface is the builtin error interface itself.
func errorInterface() *types.Interface {
	return types.Universe.Lookup("error").Type().Underlying().(*types.Interface)
}
