package files

import "fmt"

// DirNotFoundError is an error that occurs if a directory doesn't exist
type DirNotFoundError struct {
	Directory string
}

func (err DirNotFoundError) Error() string {
	return fmt.Sprintf("Directory was not found: \"%s\"", err.Directory)
}
