package a

import "errs"

// Count is how many, so no fixture signature carries a bare primitive.
type Count int

// Name is a name.
type Name string

// Flag is a yes or no.
type Flag bool

// Size is how big.
type Size int64

// good returns error last.
func good() (Count, error) { return 0, nil }

// bad returns error not last.
func bad() (error, Count) { return nil, 0 } // want "error must be the last return value"

// noErr returns no error.
func noErr() Count { return 0 }

// noResults returns nothing.
func noResults() {}

// onlyErr returns error as the sole result, which is last.
func onlyErr() error { return nil }

// named returns named results with error last.
func named() (n Count, err error) { return 0, nil }

// namedBad returns named results with error not last; the diagnostic names the
// result's own identifier, never the shared type expression.
func namedBad() (err error, n Count) { return nil, 0 } // want "error must be the last return value"

// unnamedBad returns unnamed results with error not last.
func unnamedBad() (error, Name) { return nil, "" } // want "error must be the last return value"

// Arity: the corpus must not hold the result count constant, because a bound on
// it — `len(positions) <= 3` — silences every wider signature while every
// narrower fixture reports identically.

// fourBad returns four results with the error first.
func fourBad() (error, Count, Name, Flag) { return nil, 0, "", false } // want "error must be the last return value"

// fourGood returns four results with the error last.
func fourGood() (Count, Name, Flag, error) { return 0, "", false, nil }

// fiveBad returns five results with the error second.
func fiveBad() (Count, error, Name, Flag, Size) { return 0, nil, "", false, 0 } // want "error must be the last return value"

// fiveGood returns five results with the error last.
func fiveGood() (Count, Name, Flag, Size, error) { return 0, "", false, 0, nil }

// More than one error has no compliant order, so it draws its own message rather
// than an instruction no permutation satisfies.

// twoErrs returns two errors and nothing else.
func twoErrs() (error, error) { return nil, nil } // want "more than one error has no compliant order"

// twoErrsLast returns two errors with one of them last, which does not help.
func twoErrsLast() (Count, error, error) { return 0, nil, nil } // want "more than one error has no compliant order"

// twoErrsNamed returns two named errors in ONE field; the count and the position
// are the same as twoErrsSeparate's, which is the point.
func twoErrsNamed() (clientErr, serverErr error) { return nil, nil } // want "more than one error has no compliant order"

// twoErrsSeparate returns two named errors in TWO fields.
func twoErrsSeparate() (clientErr error, serverErr error) { return nil, nil } // want "more than one error has no compliant order"

// threeErrs returns three errors, which is still one violation.
func threeErrs() (a, b, c error) { return nil, nil, nil } // want "more than one error has no compliant order"

// multiName declares two error results in one multi-name field with a compliant
// result before them; err is not last but err2 is, so this is the two-error
// violation and not the misplacement.
func multiName() (n Count, err, err2 error) { return 0, nil, nil } // want "more than one error has no compliant order"

// oneNamedErrorInAField declares one error among several names in one field, so
// exactly one result is an error and it is not last.
func oneNamedErrorInAField() (err error, a, b Count) { return nil, 0, 0 } // want "error must be the last return value"

// The sites: every function type, not only a declaration.

// Iface is an interface whose method signatures are subject to the convention.
type Iface interface {
	// Bad returns error not last.
	Bad() (error, Count) // want "error must be the last return value"
	// Good returns error last.
	Good() (Count, error)
	// OnlyErr returns error as the sole result.
	OnlyErr() error
}

// closures exercises function literals, which carry their own signatures.
func closures() {
	bad := func() (error, Count) { return nil, 0 } // want "error must be the last return value"
	good := func() (Count, error) { return 0, nil }
	_, _ = bad()
	_, _ = good()
}

// FuncField is a function-typed signature in a type definition.
type FuncField func() (error, Count) // want "error must be the last return value"

// Holder carries a function-typed struct field.
type Holder struct {
	Do func() (error, Count) // want "error must be the last return value"
}

// funcParam takes a function-typed parameter.
func funcParam(fn func() (error, Count)) { _, _ = fn() } // want "error must be the last return value"

// funcConversion converts nil to a function type written at the use site.
var funcConversion = (func() (error, Count))(nil) // want "error must be the last return value"

// funcMapValue is a map whose value type is a function type.
var funcMapValue = map[Name]func() (error, Count){} // want "error must be the last return value"

// Identity: the property is an interface implementing error, however it is
// spelled. Each of these deviates from the plain `error` case in exactly one way.

// E is an alias of the builtin error type.
type E = error

// aliasBad returns an aliased error not last.
func aliasBad() (E, Count) { return nil, 0 } // want "error must be the last return value"

// aliasGood returns an aliased error last.
func aliasGood() (Count, E) { return 0, nil }

// D is a DEFINED type over the builtin error interface, not an alias of it.
type D error

// definedBad returns a defined error type not last.
func definedBad() (D, Count) { return nil, 0 } // want "error must be the last return value"

// definedGood returns a defined error type last.
func definedGood() (Count, D) { return 0, nil }

// Wrapped is an interface EMBEDDING the builtin error interface.
type Wrapped interface{ error }

// embeddedBad returns an embedding interface not last.
func embeddedBad() (Wrapped, Count) { return nil, 0 } // want "error must be the last return value"

// embeddedGood returns an embedding interface last.
func embeddedGood() (Count, Wrapped) { return 0, nil }

// Extended is an interface embedding error and declaring more besides.
type Extended interface {
	error
	// Unwrap returns the wrapped error.
	Unwrap() error
}

// extendedBad returns a wider error interface not last.
func extendedBad() (Extended, Count) { return nil, 0 } // want "error must be the last return value"

// anonBad returns an anonymous interface with the error method, not last.
func anonBad() (interface{ Error() string }, Count) { return nil, 0 } // want "error must be the last return value"

// notAnErrorInterface has the wrong method set and is not an error.
type notAnErrorInterface interface{ Error() Name }

// notAnErrorInterfaceNotFlagged returns an interface that does not implement error.
func notAnErrorInterfaceNotFlagged() (notAnErrorInterface, Count) { return nil, 0 }

// anyNotFlagged returns the empty interface, which does not implement error.
func anyNotFlagged() (any, Count) { return nil, 0 }

// The expression form varies too: an implementation restricted to bare
// identifiers reports every fixture above identically and lets these escape.

// qualifiedAliasBad returns an imported alias of error, not last.
func qualifiedAliasBad() (errs.Alias, Count) { return nil, 0 } // want "error must be the last return value"

// qualifiedDefinedBad returns an imported defined error type, not last.
func qualifiedDefinedBad() (errs.Defined, Count) { return nil, 0 } // want "error must be the last return value"

// qualifiedEmbeddingBad returns an imported embedding interface, not last.
func qualifiedEmbeddingBad() (errs.Embedding, Count) { return nil, 0 } // want "error must be the last return value"

// qualifiedConcreteNotFlagged returns an imported CONCRETE error type, not last.
func qualifiedConcreteNotFlagged() (errs.Concrete, Count) { return errs.Concrete{}, 0 }

// starErrorBad returns a POINTER to an error interface, which is not an interface.
func starErrorNotFlagged() (*error, Count) { return nil, 0 }

// sliceErrorNotFlagged returns a slice of errors, which is not an error.
func sliceErrorNotFlagged() ([]error, Count) { return nil, 0 }

// Out of scope, each with its own reason.

// myErr is a concrete type implementing error.
type myErr struct{}

// Error makes myErr an error.
func (myErr) Error() string { return "" }

// concreteNotFlagged returns a concrete error-implementing type not last.
// Deliberately unflagged: its underlying type is not an interface, so the caller
// receives a concrete type; returning one is its own smell, not this rule.
func concreteNotFlagged() (*myErr, Count) { return nil, 0 }

// valueConcreteNotFlagged returns a concrete VALUE implementing error, not last.
func valueConcreteNotFlagged() (myErr, Count) { return myErr{}, 0 }

// genericNotFlagged returns an error-constrained type parameter not last.
// Deliberately unflagged: a type parameter's underlying type is its constraint,
// and what the caller receives is fixed at the instantiation, not here.
func genericNotFlagged[T error]() (T, Count) { var zero T; return zero, 0 }

// genericAnyNotFlagged returns an unconstrained type parameter not last.
func genericAnyNotFlagged[T any]() (T, Count) { var zero T; return zero, 0 }
