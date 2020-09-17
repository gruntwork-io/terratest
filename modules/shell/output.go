package shell

import (
	"strings"
	"sync"
)

// output contains the output after runnig a command.
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
	Lines []string
	*merged
}

func (st *outputStream) WriteString(s string) (n int, err error) {
	st.Lines = append(st.Lines, string(s))
	return st.merged.WriteString(s)
}

func (st *outputStream) String() string {
	if st == nil {
		return ""
	}

	return strings.Join(st.Lines, "\n")
}

type merged struct {
	// ensure that there are no parallel writes
	sync.Mutex
	Lines []string
}

func (m *merged) String() string {
	if m == nil {
		return ""
	}

	return strings.Join(m.Lines, "\n")
}

func (m *merged) WriteString(s string) (n int, err error) {
	m.Lock()
	defer m.Unlock()

	m.Lines = append(m.Lines, string(s))

	return len(s), nil
}
