// Command mainpkg varies the PACKAGE CLAUSE, which the rest of the corpus holds
// constant. An exemption keyed on the package name — `pass.Pkg.Name() != "main"`
// silences an entire binary — is a one-line widening inside an existing condition
// that adds no statement, so statement coverage cannot see it and every other
// fixture reports identically with it applied.
package main

// mainBad returns an error that is not last, in package main.
func mainBad() (error, Count) { return nil, 0 } // want "error must be the last return value"

// mainGood returns an error last, in package main.
func mainGood() (Count, error) { return 0, nil }

// Count is how many.
type Count int

func main() {
	_, _ = mainBad()
	_, _ = mainGood()
}
