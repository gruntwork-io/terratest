// Package parser contains methods to parse and restructure log output from go testing and terratest.
package parser

// TestResultMarker tracks the indentation level of a test result line in go test output.
// Example:
// --- FAIL: TestSnafu
//
//	--- PASS: TestSnafu/Situation
//	--- FAIL: TestSnafu/Normal
//
// The three markers for the above in order are:
// TestResultMarker{TestName: "TestSnafu", IndentLevel: 0}
// TestResultMarker{TestName: "TestSnafu/Situation", IndentLevel: 4}
// TestResultMarker{TestName: "TestSnafu/Normal", IndentLevel: 4}
type TestResultMarker struct {
	TestName    string
	IndentLevel int
}

// TestResultMarkerStack is a stack data structure to store TestResultMarkers.
type TestResultMarkerStack []TestResultMarker

// NullTestResultMarker is a blank TestResultMarker considered null. Used when peeking or popping an empty stack.
var NullTestResultMarker = TestResultMarker{}

// NULL_TEST_RESULT_MARKER is deprecated: Use NullTestResultMarker instead.
//
//nolint:staticcheck // preserving old name for backwards compatibility
var NULL_TEST_RESULT_MARKER = NullTestResultMarker //nolint:revive

// Push will push a TestResultMarker object onto the stack, returning the new one.
func (s TestResultMarkerStack) Push(v TestResultMarker) TestResultMarkerStack {
	return append(s, v)
}

// Pop will pop a TestResultMarker object off of the stack, returning the new one with the popped
// marker.
// When stack is empty, will return an empty object.
func (s TestResultMarkerStack) Pop() (TestResultMarkerStack, TestResultMarker) {
	l := len(s)
	if l == 0 {
		return s, NullTestResultMarker
	}

	return s[:l-1], s[l-1]
}

// Peek will return the top TestResultMarker from the stack, but will not remove it.
func (s TestResultMarkerStack) Peek() TestResultMarker {
	l := len(s)
	if l == 0 {
		return NullTestResultMarker
	}

	return s[l-1]
}

// IsEmpty will return whether or not the stack is empty.
func (s TestResultMarkerStack) IsEmpty() bool {
	return len(s) == 0
}

// RemoveDedentedTestResultMarkers will pop items off of the stack of TestResultMarker objects until the top most item
// has an indent level less than the current indent level.
// Assumes that the stack is ordered, in that recently pushed items in the stack have higher indent levels.
func (s TestResultMarkerStack) RemoveDedentedTestResultMarkers(currentIndentLevel int) TestResultMarkerStack {
	// This loop is a garbage collection of the stack, where it removes entries every time we dedent out of a fail
	// block.
	for !s.IsEmpty() && s.Peek().IndentLevel >= currentIndentLevel {
		s, _ = s.Pop()
	}

	return s
}
