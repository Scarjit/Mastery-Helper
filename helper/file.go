package helper

import (
	"errors"
	"os"
)

// FileExists checks if a file exists
func FileExists(path string) bool {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return false
	}
	return true
}
