package errlast_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/tools/go/analysis/analysistest"

	errlast "github.com/gomatic/yze-errlast"
)

func TestErrorMustBeLast(t *testing.T) {
	analysistest.Run(t, analysistest.TestData(), errlast.Analyzer, "a")
}

func TestRegistrationIsWellFormed(t *testing.T) {
	assert.NoError(t, errlast.Registration.Validate())
	assert.Equal(t, "yze/errlast", errlast.Registration.RuleID())
	assert.Same(t, errlast.Analyzer, errlast.Registration.Analyzer)
}
