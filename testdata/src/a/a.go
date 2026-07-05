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

// E is an alias of the builtin error type; results typed through it must obey the convention.
type E = error

// aliasBad returns an aliased error not last.
func aliasBad() (E, int) { return nil, 0 } // want "error must be the last return value"

// aliasGood returns an aliased error last.
func aliasGood() (int, E) { return 0, nil }

// myErr is a concrete type implementing error.
type myErr struct{}

func (myErr) Error() string { return "" }

// concreteNotFlagged returns a concrete error-implementing type not last.
// Deliberately unflagged: the analyzer checks only the builtin error interface
// type; returning a concrete error type is its own smell, not this rule.
func concreteNotFlagged() (*myErr, int) { return nil, 0 }

// genericNotFlagged returns an error-constrained type parameter not last.
// Deliberately unflagged: a type parameter is not the builtin error interface
// type, even when constrained by it.
func genericNotFlagged[T error]() (T, int) { var zero T; return zero, 0 }

// multiName declares two error results in one multi-name field; flattened, err
// sits at position 1 of 3, which is not last, so the shared type expression is
// flagged once (err2, at the final position, is compliant).
func multiName() (n int, err, err2 error) { return 0, nil, nil } // want "error must be the last return value"
