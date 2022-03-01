package testdata

import (
	"path/filepath"
	"runtime"
)

// Dir returns the directory of the testdata module.
func Dir() string {
	_, currentFile, _, _ := runtime.Caller(0)
	return filepath.Dir(currentFile)
}

// Filename returns the file in testdata module.
func Filename(filename string) string {
	return filepath.Join(Dir(), filename)
}
