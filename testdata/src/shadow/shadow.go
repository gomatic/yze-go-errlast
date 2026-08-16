// Package shadow declares a package-scope type named error, shadowing the
// predeclared identifier. It is its own package because the declaration changes
// what `error` means for every other file beside it.
//
// It exists because the revision this fixture was written against recognised an
// error result by its OBJECT — a *types.Named at universe scope literally named
// "error" — and was silent here, where the result is an ordinary interface with
// the error method and a value of it satisfies the builtin error. Nothing in the
// repository could fail on removing that guard, and removing it survived the full
// gate at 100.0% coverage. Recognising the PROPERTY rather than the object makes
// this shape reportable and the guard unnecessary.
package shadow

// error shadows the predeclared identifier with an ordinary interface.
type error interface{ Error() string }

// shadowedBad returns a shadowing error interface that is not last.
func shadowedBad() (error, Count) { return nil, 0 } // want "error must be the last return value"

// shadowedGood returns a shadowing error interface last.
func shadowedGood() (Count, error) { return 0, nil }

// Count is how many.
type Count int

var _, _ = shadowedBad, shadowedGood
