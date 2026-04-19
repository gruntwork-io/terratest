package aws

// SleepWithContextForTest exposes the unexported sleepWithContext helper to the external
// aws_test package. It exists only in tests (file has the _test.go suffix) and is not part of
// the public API.
var SleepWithContextForTest = sleepWithContext
