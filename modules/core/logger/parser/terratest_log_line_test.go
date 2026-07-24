package parser_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/core/v2/logger/parser"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIsIndentedTerratestLogLine(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string
		in   string
		out  bool
	}{
		{
			name: "IndentedTerratestLine",
			in:   "    apply_test.go:42: TestFoo 2026-07-18T13:36:46-04:00 logger.go:81: applying",
			out:  true,
		},
		{
			name: "IndentedSubtestLine",
			in:   "        apply_test.go:42: TestFoo/Sub1 2026-07-18T13:36:46-04:00 logger.go:81: applying",
			out:  true,
		},
		{
			name: "UnindentedTerratestLine",
			in:   "TestFoo 2026-07-18T13:36:46-04:00 logger.go:81: applying",
			out:  false,
		},
		{
			name: "PlainTLogLine",
			in:   "    apply_test.go:42: some plain message",
			out:  false,
		},
		{
			name: "IndentedResultLine",
			in:   "    --- PASS: TestFoo (0.02s)",
			out:  false,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, testCase.out, parser.IsIndentedTerratestLogLine(testCase.in))
		})
	}
}

func TestGetTestNameFromIndentedTerratestLogLine(t *testing.T) {
	t.Parallel()

	assert.Equal(t, "TestFoo", parser.GetTestNameFromIndentedTerratestLogLine(
		"    apply_test.go:42: TestFoo 2026-07-18T13:36:46-04:00 logger.go:81: applying"))
	assert.Equal(t, "TestFoo/Sub1", parser.GetTestNameFromIndentedTerratestLogLine(
		"        apply_test.go:42: TestFoo/Sub1 2026-07-18T13:36:46-04:00 logger.go:81: applying"))
}

// TestSpawnParsersDeinterleavesTLogOutput verifies that the parser de-interleaves parallel-test output that terratest
// now emits through t.Log (indented and decorated by the testing framework). Both tests resume up front, so the
// `=== CONT` status lines cannot distinguish them; correct attribution relies on the test name embedded in each line.
func TestSpawnParsersDeinterleavesTLogOutput(t *testing.T) {
	t.Parallel()

	// Interleaved `go test -v` output (non-JSON) as produced after logging is routed through t.Log.
	sample := strings.Join([]string{
		"=== RUN   TestParA",
		"=== PAUSE TestParA",
		"=== RUN   TestParB",
		"=== PAUSE TestParB",
		"=== CONT  TestParA",
		"=== CONT  TestParB",
		"    a_test.go:10: TestParA 2026-07-18T13:36:46-04:00 logger.go:81: MARKA payload 0",
		"    b_test.go:20: TestParB 2026-07-18T13:36:46-04:00 logger.go:81: MARKB payload 0",
		"    b_test.go:20: TestParB 2026-07-18T13:36:46-04:00 logger.go:81: MARKB payload 1",
		"    a_test.go:10: TestParA 2026-07-18T13:36:46-04:00 logger.go:81: MARKA payload 1",
		"    a_test.go:10: TestParA 2026-07-18T13:36:46-04:00 logger.go:81: MARKA payload 2",
		"    b_test.go:20: TestParB 2026-07-18T13:36:46-04:00 logger.go:81: MARKB payload 2",
		"--- PASS: TestParA (0.02s)",
		"--- PASS: TestParB (0.02s)",
		"PASS",
		"ok  \tpkg\t0.10s",
		"",
	}, "\n")

	out := t.TempDir()
	parser.SpawnParsers(logrus.New(), strings.NewReader(sample), out)

	readLog := func(test string) string {
		b, err := os.ReadFile(filepath.Join(out, test+".log"))
		require.NoError(t, err, "expected a log file for %s", test)

		return string(b)
	}

	a := readLog("TestParA")
	assert.Equal(t, 3, strings.Count(a, "MARKA"), "TestParA.log should contain all of its own lines")
	assert.NotContains(t, a, "MARKB", "TestParA.log should not contain TestParB output")

	b := readLog("TestParB")
	assert.Equal(t, 3, strings.Count(b, "MARKB"), "TestParB.log should contain all of its own lines")
	assert.NotContains(t, b, "MARKA", "TestParB.log should not contain TestParA output")
}
