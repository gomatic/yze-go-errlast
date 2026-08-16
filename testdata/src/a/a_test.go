package a

// This file varies the FILE KIND, which the rest of the corpus holds constant at
// non-test. A suffix test on the reported file's name — silencing every finding
// in a _test.go file — is a one-line widening inside an existing condition that
// adds no statement, so statement coverage cannot see it and every other fixture
// reports identically with it applied. errlast is absent from the runner's
// sourceOnly map, so it reports in test files under the runner as well.

// testFileBad returns an error that is not last, in a test file.
func testFileBad() (error, Count) { return nil, 0 } // want "error must be the last return value"

// testFileGood returns an error last, in a test file.
func testFileGood() (Count, error) { return 0, nil }

var _, _ = testFileBad, testFileGood
