package utils

import (
	"go/format"
	"io/ioutil"
	"os"
	"path/filepath"

	"golang.org/x/tools/imports"
)

// GoFormat formats go file in canonical gofmt style and fix import statements.
func GoFormat(file string) error {
	c, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	c, err = imports.Process(file, c, nil)
	if err != nil {
		return err
	}
	c, err = format.Source(c)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(file, c, 0666)
}

func FileExists(file string) bool {
	if _, err := os.Stat(file); err != nil {
		return false
	}
	return true
}

func DeleteFileByPattern(pattern string) error {
	files, err := filepath.Glob(pattern)
	if err != nil {
		return err
	}
	for _, file := range files {
		err = os.Remove(file)
		if err != nil {
			return err
		}
	}
	return nil
}
