package parser_test

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/core/v2/logger/parser"
	"github.com/stretchr/testify/assert"
)

func createTestStack() parser.TestResultMarkerStack {
	return parser.TestResultMarkerStack{
		{
			TestName:    "TestSnafu",
			IndentLevel: 0,
		},
		{
			TestName:    "TestSnafu/Situation",
			IndentLevel: 4,
		},
		{
			TestName:    "TestSnafu/Normal",
			IndentLevel: 4,
		},
	}
}

// TestStackPush tests that pushing items to the stack appends to the list.
func TestStackPush(t *testing.T) {
	t.Parallel()

	markers := createTestStack()
	newMarker := parser.TestResultMarker{
		TestName:    "TestThatEverythingWorks",
		IndentLevel: 0,
	}
	markers = markers.Push(newMarker)
	assert.Len(t, markers, 4)
	assert.Equal(t, newMarker, markers[3])
}

// TestStackPop tests that popping items off the stack will remove and return the LAST item.
func TestStackPop(t *testing.T) {
	t.Parallel()

	originalMarkers := createTestStack()
	markers := createTestStack()

	markers, poppedMarker := markers.Pop()
	assert.Equal(t, originalMarkers[2], poppedMarker)
	assert.Len(t, markers, 2)

	markers, poppedMarker = markers.Pop()
	assert.Equal(t, originalMarkers[1], poppedMarker)
	assert.Len(t, markers, 1)

	markers, poppedMarker = markers.Pop()
	assert.Equal(t, originalMarkers[0], poppedMarker)
	assert.Empty(t, markers)
}

// TestStackPopEmpty tests that popping item off an empty stack returns an empty TestResultMarker.
func TestStackPopEmpty(t *testing.T) {
	t.Parallel()

	markers := parser.TestResultMarkerStack{}

	markers, poppedMarker := markers.Pop()
	assert.Empty(t, markers)
	assert.Equal(t, parser.NullTestResultMarker, poppedMarker)
}

// TestPeek tests that peek returns the LAST item in the list WITHOUT removing it.
func TestPeek(t *testing.T) {
	t.Parallel()

	originalMarkers := createTestStack()
	markers := createTestStack()
	peekedMarker := markers.Peek()
	assert.Equal(t, originalMarkers[2], peekedMarker)
	assert.Equal(t, originalMarkers, markers)
}

// TestPeekEmpty tests that peeking an empty stack returns an empty TestResultMarker.
func TestPeekEmpty(t *testing.T) {
	t.Parallel()

	markers := parser.TestResultMarkerStack{}
	peekedMarker := markers.Peek()
	assert.Empty(t, markers)
	assert.Equal(t, parser.NullTestResultMarker, peekedMarker)
}

// TestIsEmpty tests that IsEmpty only returns true on empty stack.
func TestIsEmpty(t *testing.T) {
	t.Parallel()

	emptyMarkerStack := parser.TestResultMarkerStack{}
	fullMarkerStack := createTestStack()

	assert.True(t, emptyMarkerStack.IsEmpty())
	assert.False(t, fullMarkerStack.IsEmpty())
}

// TestRemoveDedentedTestResultMarkers tests that items dedented from the current level are removed.
func TestRemoveDedentedTestResultMarkers(t *testing.T) {
	t.Parallel()

	originalMarkers := createTestStack()
	newMarkers := originalMarkers.RemoveDedentedTestResultMarkers(2)
	assert.Len(t, newMarkers, 1)
	assert.Equal(t, newMarkers, originalMarkers[:1])
}

// TestRemoveDedentedTestResultMarkersEmpty tests that RemoveDedentedTestResultMarkers handles empty stack.
func TestRemoveDedentedTestResultMarkersEmpty(t *testing.T) {
	t.Parallel()

	originalMarkers := parser.TestResultMarkerStack{}
	newMarkers := originalMarkers.RemoveDedentedTestResultMarkers(2)
	assert.Empty(t, newMarkers)
}

// TestRemoveDedentedTestResultMarkersAll tests that RemoveDedentedTestResultMarkers handles removing everything.
func TestRemoveDedentedTestResultMarkersAll(t *testing.T) {
	t.Parallel()

	originalMarkers := parser.TestResultMarkerStack{}
	newMarkers := originalMarkers.RemoveDedentedTestResultMarkers(-1)
	assert.Empty(t, newMarkers)
}
