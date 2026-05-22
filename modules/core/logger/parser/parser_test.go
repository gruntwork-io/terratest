package parser_test

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/core/v2/logger/parser"
	"github.com/stretchr/testify/assert"
)

func TestGetIndent(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string
		in   string
		out  string
	}{
		{
			name: "BaseCase",
			in:   "    --- FAIL: TestSnafu",
			out:  "    ",
		},
		{
			name: "NoIndent",
			in:   "--- FAIL: TestSnafu",
			out:  "",
		},
		{
			name: "EmptyString",
			in:   "",
			out:  "",
		},
		{
			name: "Tabs",
			in:   "\t\t---FAIL: TestSnafu",
			out:  "\t\t",
		},
		{
			name: "MixTabSpace",
			in:   "\t    ---FAIL: TestSnafu",
			out:  "\t    ",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(
				t,
				testCase.out,
				parser.GetIndent(testCase.in),
			)
		})
	}
}

func TestGetTestNameFromResultLine(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string
		in   string
		out  string
	}{
		{
			name: "BaseCase",
			in:   "--- PASS: TestGetTestNameFromResultLine (0.00s)",
			out:  "TestGetTestNameFromResultLine",
		},
		{
			name: "Indented",
			in:   "    --- PASS: TestGetTestNameFromResultLine/Indented (0.00s)",
			out:  "TestGetTestNameFromResultLine/Indented",
		},
		{
			name: "SpecialChars",
			in:   "    --- PASS: TestGetTestNameFromResultLine/SpecialChars---_FAIL (0.00s)",
			out:  "TestGetTestNameFromResultLine/SpecialChars---_FAIL",
		},
		{
			name: "WhenFailed",
			in:   "--- FAIL: TestGetTestNameFromResultLine (0.00s)",
			out:  "TestGetTestNameFromResultLine",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(
				t,
				testCase.out,
				parser.GetTestNameFromResultLine(testCase.in),
			)
		})
	}
}

func TestIsResultLine(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string
		in   string
		out  bool
	}{
		{
			name: "BaseCase",
			in:   "--- PASS: TestIsResultLine (0.00s)",
			out:  true,
		},
		{
			name: "Indented",
			in:   "    --- PASS: TestIsResultLine/Indented (0.00s)",
			out:  true,
		},
		{
			name: "SpecialChars",
			in:   "    --- PASS: TestIsResultLine/SpecialChars---_FAIL (0.00s)",
			out:  true,
		},
		{
			name: "WhenFailed",
			in:   "--- FAIL: TestIsResultLine (0.00s)",
			out:  true,
		},
		{
			name: "NonResultLine",
			in:   "=== RUN TestIsResultLine",
			out:  false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(
				t,
				testCase.out,
				parser.IsResultLine(testCase.in),
			)
		})
	}
}

func TestGetTestNameFromStatusLine(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string
		in   string
		out  string
	}{
		{
			name: "BaseCase",
			in:   "=== RUN   TestGetTestNameFromStatusLine",
			out:  "TestGetTestNameFromStatusLine",
		},
		{
			name: "Indented",
			in:   "    === RUN   TestGetTestNameFromStatusLine/Indented",
			out:  "TestGetTestNameFromStatusLine/Indented",
		},
		{
			name: "SpecialChars",
			in:   "=== RUN   TestGetTestNameFromStatusLine/SpecialChars---_FAIL",
			out:  "TestGetTestNameFromStatusLine/SpecialChars---_FAIL",
		},
		{
			name: "WhenPaused",
			in:   "=== PAUSE TestGetTestNameFromStatusLine",
			out:  "TestGetTestNameFromStatusLine",
		},
		{
			name: "WhenCont",
			in:   "=== CONT  TestGetTestNameFromStatusLine",
			out:  "TestGetTestNameFromStatusLine",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(
				t,
				testCase.out,
				parser.GetTestNameFromStatusLine(testCase.in),
			)
		})
	}
}

func TestIsStatusLine(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string
		in   string
		out  bool
	}{
		{
			name: "BaseCase",
			in:   "=== RUN   TestGetTestNameFromStatusLine",
			out:  true,
		},
		{
			name: "Indented",
			in:   "    === RUN   TestGetTestNameFromStatusLine/Indented",
			out:  true,
		},
		{
			name: "SpecialChars",
			in:   "=== RUN   TestGetTestNameFromStatusLine/SpecialChars---_FAIL",
			out:  true,
		},
		{
			name: "WhenPaused",
			in:   "=== PAUSE TestGetTestNameFromStatusLine",
			out:  true,
		},
		{
			name: "WhenCont",
			in:   "=== CONT  TestGetTestNameFromStatusLine",
			out:  true,
		},
		{
			name: "NonStatusLine",
			in:   "--- FAIL: TestIsStatusLine",
			out:  false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(
				t,
				testCase.out,
				parser.IsStatusLine(testCase.in),
			)
		})
	}
}

func TestIsSummaryLine(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string
		in   string
		out  bool
	}{
		{
			name: "BaseCase",
			in:   "ok  	github.com/gruntwork-io/terratest/test	812.034s",
			out:  true,
		},
		{
			name: "NotSummary",
			in:   "--- FAIL: TestIsStatusLine",
			out:  false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(
				t,
				testCase.out,
				parser.IsSummaryLine(testCase.in),
			)
		})
	}
}

func TestIsPanicLine(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string
		in   string
		out  bool
	}{
		{
			name: "BaseCase",
			in:   "panic: error [recovered]",
			out:  true,
		},
		{
			name: "NotPanic",
			in:   "--- FAIL: TestIsStatusLine",
			out:  false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(
				t,
				testCase.out,
				parser.IsPanicLine(testCase.in),
			)
		})
	}
}
