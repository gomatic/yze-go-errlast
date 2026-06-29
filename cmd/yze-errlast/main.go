// Command yze-errlast runs the errlast analyzer as a standalone go/analysis
// checker (text and -json output, and as a `go vet -vettool`).
package main

import (
	"golang.org/x/tools/go/analysis/singlechecker"

	errlast "github.com/gomatic/yze-errlast"
)

// run is the analysis entry point, indirected so the binary's wiring is testable
// without invoking the real driver (which loads packages and exits the process).
var run = singlechecker.Main

func main() { run(errlast.Analyzer) }
