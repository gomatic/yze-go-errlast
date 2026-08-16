package errlast

import (
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

// The corpus Finding is {rule, path, line} and analysistest matches `want`
// directives by LINE, so neither can see a column and neither can distinguish two
// diagnostics that share a position from one. Both are part of the standard: the
// position is what an editor opens and what a quota charges. They are asserted
// here, against a pass built by hand, and written from the package doc comment
// rather than from a run — moving the report from the result to its shared type
// expression shifts every column and survives the whole corpus in silence.

// source is a complete Go file a case type-checks and analyzes.
type source string

// report is one diagnostic, projected to the parts the corpus cannot carry.
type report struct {
	message string
	line    int
	column  int
}

// analyze type-checks the source and runs the analyzer over it, collecting every
// diagnostic. A constructed pass is used rather than a fixture because a fixture
// cannot assert a column.
func analyze(t *testing.T, src source) []report {
	t.Helper()
	return analyzeWith(t, src, recordTypes)
}

// typeRecording is what a pass's TypesInfo carries into the analyzer. It is a
// parameter because an analyzer does not get to assume the driver recorded a type
// for every expression, and the refusal for the case where it did not is otherwise
// unreachable — see TestAResultWithNoRecordedTypeIsRefused.
type typeRecording bool

const (
	// recordTypes is the ordinary pass: the driver recorded every expression.
	recordTypes typeRecording = true
	// dropTypes is a pass whose TypesInfo has nothing in it, so TypeOf returns nil.
	dropTypes typeRecording = false
)

func analyzeWith(t *testing.T, src source, recording typeRecording) []report {
	t.Helper()
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "p.go", string(src), parser.SkipObjectResolution)
	require.NoError(t, err)
	info := &types.Info{
		Types: map[ast.Expr]types.TypeAndValue{},
		Defs:  map[*ast.Ident]types.Object{},
		Uses:  map[*ast.Ident]types.Object{},
	}
	files := []*ast.File{file}
	pkg, err := (&types.Config{}).Check("p", fset, files, info)
	require.NoError(t, err)
	if recording == dropTypes {
		info = &types.Info{}
	}
	var got []report
	_, err = run(&analysis.Pass{
		Fset:      fset,
		Files:     files,
		Pkg:       pkg,
		TypesInfo: info,
		ResultOf:  map[*analysis.Analyzer]any{inspect.Analyzer: inspector.New(files)},
		Report: func(diagnostic analysis.Diagnostic) {
			position := fset.Position(diagnostic.Pos)
			got = append(got, report{diagnostic.Message, position.Line, position.Column})
		},
	})
	require.NoError(t, err)
	return got
}

// TestUnnamedResultIsReportedAtItsTypeExpression pins the column for the shape
// that has no identifier to name: the type expression is the only position the
// result has.
func TestUnnamedResultIsReportedAtItsTypeExpression(t *testing.T) {
	got := analyze(t, "package p\n\nfunc f() (error, int) { return nil, 0 }\n")
	assert.Equal(t, []report{{line: 3, column: 11, message: message}}, got)
}

// TestNamedResultIsReportedAtItsOwnIdentifier pins the column for a named result.
// It is the identifier at column 11, NOT the type expression at column 15, because
// two results sharing one type expression must not share one position.
func TestNamedResultIsReportedAtItsOwnIdentifier(t *testing.T) {
	got := analyze(t, "package p\n\nfunc f() (err error, n int) { return nil, 0 }\n")
	assert.Equal(t, []report{{line: 3, column: 11, message: message}}, got)
}

// TestMultiNameFieldReportsOnceAtTheFirstError is the case an earlier revision
// got wrong in two ways at once: it emitted one diagnostic per NAME, all of them
// at the shared type expression, so -json carried two entries at an identical
// position while the text driver deduplicated them out of sight.
func TestMultiNameFieldReportsOnceAtTheFirstError(t *testing.T) {
	got := analyze(t, "package p\n\nfunc f() (err, err2 error, n int) { return nil, nil, 0 }\n")
	assert.Equal(t, []report{{line: 3, column: 11, message: multiErrorMessage}}, got)
}

// TestTheCountDoesNotFollowTheSpelling is the contract the previous test's
// position serves: one field of two names and two fields of one name are the same
// violation, so they draw the same single diagnostic at the same result.
func TestTheCountDoesNotFollowTheSpelling(t *testing.T) {
	oneField := analyze(t, "package p\n\nfunc f() (err, err2 error) { return nil, nil }\n")
	twoFields := analyze(t, "package p\n\nfunc f() (err error, err2 error) { return nil, nil }\n")
	assert.Equal(t, []report{{line: 3, column: 11, message: multiErrorMessage}}, oneField)
	assert.Equal(t, oneField, twoFields)
}

// TestEveryResultIsCountedEvenWhenItShareAField guards the arithmetic the index
// comparison rests on: the error is last only because both preceding names count.
func TestEveryResultIsCountedEvenWhenItSharesAField(t *testing.T) {
	assert.Empty(t, analyze(t, "package p\n\nfunc f() (a, b int, err error) { return 0, 0, nil }\n"))
}

// TestSingleErrorLastIsSilent pins the other side of the index boundary.
func TestSingleErrorLastIsSilent(t *testing.T) {
	assert.Empty(t, analyze(t, "package p\n\nfunc f() (n int, err error) { return 0, nil }\n"))
}

// TestMultipleErrorsAreOneViolationEvenWhenOneIsLast pins the boundary between the
// two messages: the count of errors decides, not their placement, because placement
// cannot fix it.
func TestMultipleErrorsAreOneViolationEvenWhenOneIsLast(t *testing.T) {
	got := analyze(t, "package p\n\nfunc f() (n int, a error, b error) { return 0, nil, nil }\n")
	assert.Equal(t, []report{{line: 3, column: 18, message: multiErrorMessage}}, got)
}

// TestAResultWithNoRecordedTypeIsRefused pins the nil arm of isErrorResult, which
// no fixture can reach: a corpus package always type-checks, so TypeOf always
// answers. The arm shares its `return false` statement with the type-parameter
// arm, so statement coverage reports it green while nothing enters it — and
// deleting `nil` from the label is not cosmetic, because a nil types.Type falls
// through to isErrorInterface, which calls Underlying on it and segfaults. The
// contract is that an analyzer does not assume its driver recorded every
// expression, so a result with no recorded type is refused rather than fatal.
func TestAResultWithNoRecordedTypeIsRefused(t *testing.T) {
	src := source("package p\n\nfunc f() (error, int) { return nil, 0 }\n")
	require.NotEmpty(t, analyzeWith(t, src, recordTypes), "the same source must report when types ARE recorded")
	assert.Empty(t, analyzeWith(t, src, dropTypes))
}

// TestMessagesAreTheDocumentedOnes pins the text an author reads. The single-error
// message prescribes a move, which is always available; the multi-error message
// must NOT prescribe an order, because no order exists.
func TestMessagesAreTheDocumentedOnes(t *testing.T) {
	assert.Equal(t, "error must be the last return value", message)
	assert.Equal(
		t,
		"a signature returning more than one error has no compliant order; return one joined error, or one value carrying them",
		multiErrorMessage,
	)
	assert.NotContains(t, multiErrorMessage, "last return value")
}
