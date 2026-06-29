package a

// good returns error last.
func good() (int, error) { return 0, nil }

// bad returns error not last.
func bad() (error, int) { return nil, 0 } // want "error must be the last return value"

// noErr returns no error.
func noErr() int { return 0 }

// noResults returns nothing.
func noResults() {}

// onlyErr returns error as the sole result, which is last.
func onlyErr() error { return nil }

// twoErrs returns two errors; the first is not last.
func twoErrs() (error, error) { return nil, nil } // want "error must be the last return value"

// named returns named results with error last.
func named() (n int, err error) { return 0, nil }

// unnamedBad returns unnamed results with error not last.
func unnamedBad() (error, string) { return nil, "" } // want "error must be the last return value"

// Iface is an interface whose method signatures are subject to the convention.
type Iface interface {
	// Bad returns error not last.
	Bad() (error, int) // want "error must be the last return value"
	// Good returns error last.
	Good() (int, error)
	// OnlyErr returns error as the sole result.
	OnlyErr() error
}

// closures exercises function literals, which carry their own signatures.
func closures() {
	bad := func() (error, int) { return nil, 0 } // want "error must be the last return value"
	good := func() (int, error) { return 0, nil }
	_, _ = bad()
	_, _ = good()
}

// FuncField is a function-typed signature in a type definition.
type FuncField func() (error, int) // want "error must be the last return value"
