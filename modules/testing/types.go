package testing

// TestingT is an interface that describes the implementation of the testing object
// that the majority of Terratest functions accept as first argument.
// Using an interface that describes testing.T instead of the actual implementation
// makes terratest usable in a wider variety of contexts (e.g. use with ginkgo : https://godoc.org/github.com/onsi/ginkgo#GinkgoT)
type TestingT interface {
	//Fail marks the function as having failed but continues execution.
	Fail()
	// FailNow marks the function as having failed and stops its execution
	// by calling runtime.Goexit (which then runs all deferred calls in the
	// current goroutine).
	// Execution will continue at the next test or benchmark.
	// FailNow must be called from the goroutine running the
	// test or benchmark function, not from other goroutines
	// created during the test. Calling FailNow does not stop
	// those other goroutines.
	FailNow()
	Fatal(args ...interface{})
	// Fatalf is equivalent to Logf followed by FailNow.
	Fatalf(format string, args ...interface{})
	// Error is equivalent to Log followed by Fail.
	Error(args ...interface{})
	// Errorf is equivalent to Logf followed by Fail.
	Errorf(format string, args ...interface{})
	// Name returns the name of the running test or benchmark.
	Name() string
}
