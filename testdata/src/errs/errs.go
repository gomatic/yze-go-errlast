// Package errs exists so a fixture can name an error type through an IMPORT.
// Without it the corpus never varies the expression form of a result type: every
// result is a bare identifier, and an implementation restricted to bare
// identifiers reports identically on the whole corpus while a qualified alias
// escapes it silently.
package errs

// Alias is the builtin error type under another package's name.
type Alias = error

// Defined is a defined type over the builtin error interface.
type Defined error

// Embedding is an interface embedding the builtin error interface.
type Embedding interface{ error }

// Concrete implements error without being an interface.
type Concrete struct{}

// Error makes Concrete an error.
func (Concrete) Error() string { return "" }
