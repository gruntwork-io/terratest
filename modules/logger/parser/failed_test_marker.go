// Package logger/parser contains methods to parse and restructure log output from go testing and terratest
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

// TestResultMarkerStack is a stack data structure to store TestResultMarkers
type TestResultMarkerStack []TestResultMarker

// A blank TestResultMarker is considered null. Used when peeking or popping an empty stack.
var NULL_TEST_RESULT_MARKER = TestResultMarker{}

// TestResultMarker.push will push a TestResultMarker object onto the stack, returning the new one.
func (s TestResultMarkerStack) push(v TestResultMarker) TestResultMarkerStack {
	return append(s, v)
}

// TestResultMarker.pop will pop a TestResultMarker object off of the stack, returning the new one with the popped
// marker.
// When stack is empty, will return an empty object.
func (s TestResultMarkerStack) pop() (TestResultMarkerStack, TestResultMarker) {
	l := len(s)
	if l == 0 {
		return s, NULL_TEST_RESULT_MARKER
	}
	return s[:l-1], s[l-1]
}

// TestResultMarker.peek will return the top TestResultMarker from the stack, but will not remove it.
func (s TestResultMarkerStack) peek() TestResultMarker {
	l := len(s)
	if l == 0 {
		return NULL_TEST_RESULT_MARKER
	}
	return s[l-1]
}

// TestResultMarker.isEmpty will return whether or not the stack is empty.
func (s TestResultMarkerStack) isEmpty() bool {
	return len(s) == 0
}

// removeDedentedTestResultMarkers will pop items off of the stack of TestResultMarker objects until the top most item
// has an indent level less than the current indent level.
// Assumes that the stack is ordered, in that recently pushed items in the stack have higher indent levels.
func (s TestResultMarkerStack) removeDedentedTestResultMarkers(currentIndentLevel int) TestResultMarkerStack {
	// This loop is a garbage collection of the stack, where it removes entries every time we dedent out of a fail
	// block.
	for !s.isEmpty() && s.peek().IndentLevel >= currentIndentLevel {
		s, _ = s.pop()
	}
	return s
}
