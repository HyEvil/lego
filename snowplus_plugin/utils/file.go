package utils

import (
	"github.com/kubernetes/gengo/parser"
	"go/format"
	"io/ioutil"
	"os"
	"strings"

	"golang.org/x/tools/imports"
)

// OpenOrCreate opens a file if it exists, otherwise creates it.
// If the file path contains directories, it will make them first.
func OpenOrCreate(file string) (*os.File, error) {
	if FileExists(file) {
		return os.OpenFile(file, os.O_RDWR|os.O_APPEND, 0644)
	}
	if i := strings.LastIndex(file, "/"); i != -1 {
		if err := os.MkdirAll(file[:i], 0755); err != nil {
			return nil, err
		}
	}
	return os.Create(file)
}

// Creates or truncates the named file,if the file path contains directories, it will make them first.
func Create(file string) (*os.File, error) {
	if !FileExists(file) {
		if i := strings.LastIndex(file, "/"); i != -1 {
			if err := os.MkdirAll(file[:i], 0755); err != nil {
				return nil, err
			}
		}
	}
	return os.Create(file)
}

// FileExists checks whether a file exists.
func FileExists(file string) bool {
	parser.New().AddDir("")
	if _, err := os.Stat(file); os.IsNotExist(err) {
		return false
	}
	return true
}

// GoFormat formats go file in canonical gofmt style and fix import statements.
func GoFormat(file string) error {
	c, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	c, err = format.Source(c)
	if err != nil {
		return err
	}
	c, err = imports.Process(file, c, nil)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(file, c, 0666)
}

func FormatStr(str string, path string) (string, error) {
	data, err := format.Source([]byte(str))
	if err != nil {
		return "", err
	}
	data, err = imports.Process(path, data, nil)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func IsDir(path string) bool {
	if file, err := os.Stat(path); err == nil {
		return file.IsDir()
	}
	return false
}
