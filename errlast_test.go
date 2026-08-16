package errlast_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/tools/go/analysis/analysistest"

	errlast "github.com/gomatic/yze-go-errlast"
)

// TestErrorMustBeLast runs the corpus. The packages are separate because each
// varies a dimension the others hold constant: `a` carries the rule and its
// boundaries plus a _test.go sibling for the file kind, `mainpkg` varies the
// package clause, `shadow` varies whether `error` is the predeclared identifier,
// and `a` imports `errs` so a result type is spelled as a qualified expression
// rather than a bare identifier. A widening keyed on any of them reports
// identically on a corpus that holds it constant.
func TestErrorMustBeLast(t *testing.T) {
	analysistest.Run(t, analysistest.TestData(), errlast.Analyzer, "a", "mainpkg", "shadow")
}

func TestRegistrationIsWellFormed(t *testing.T) {
	assert.NoError(t, errlast.Registration.Validate())
	assert.Equal(t, "yze/errlast", errlast.Registration.RuleID())
	assert.Same(t, errlast.Analyzer, errlast.Registration.Analyzer)
}
