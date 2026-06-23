package shell

import (
	"strings"
	"sync"
)

// output contains the output after running a command.
type output struct {
	stdout *outputStream
	stderr *outputStream
	// merged contains stdout  and stderr merged into one stream.
	merged *merged
}

func newOutput() *output {
	m := new(merged)

	return &output{
		merged: m,
		stdout: &outputStream{
			merged: m,
		},
		stderr: &outputStream{
			merged: m,
		},
	}
}

func (o *output) Stdout() string {
	if o == nil {
		return ""
	}

	return o.stdout.String()
}

func (o *output) Stderr() string {
	if o == nil {
		return ""
	}

	return o.stderr.String()
}

func (o *output) Combined() string {
	if o == nil {
		return ""
	}

	return o.merged.String()
}

type outputStream struct {
	merged *merged
	lines  []string
}

func (st *outputStream) WriteString(s string) (n int, err error) {
	st.lines = append(st.lines, s)

	return st.merged.WriteString(s)
}

func (st *outputStream) String() string {
	if st == nil {
		return ""
	}

	return strings.Join(st.lines, "\n")
}

type merged struct {
	lines []string
	// ensure that there are no parallel writes
	sync.Mutex
}

func (m *merged) String() string {
	if m == nil {
		return ""
	}

	return strings.Join(m.lines, "\n")
}

func (m *merged) WriteString(s string) (n int, err error) {
	m.Lock()
	defer m.Unlock()

	m.lines = append(m.lines, s)

	return len(s), nil
}
