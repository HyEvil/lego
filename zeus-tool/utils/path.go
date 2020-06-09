package utils

import (
	"os"
	"path/filepath"
)

func ProjectPath(path ...string) string {
	curPath, err := os.Getwd()
	if err != nil {
		return ""
	}
	paths := []string{curPath}
	paths = append(paths, path...)
	return filepath.Join(paths...)
}

func ProjectDirPathEx(create bool, path ...string) string {
	curPath, err := os.Getwd()
	if err != nil {
		return ""
	}
	paths := []string{curPath}
	paths = append(paths, path...)
	p := filepath.Join(paths...)
	if create {
		os.Mkdir(p, os.ModePerm)
	}
	return p
}

func ProjectFilePathEx(create bool, path ...string) string {
	curPath, err := os.Getwd()
	if err != nil {
		return ""
	}
	paths := []string{curPath}
	paths = append(paths, path...)
	p := filepath.Join(paths...)
	if create {
		if !FileExists(filepath.Dir(p)){
			os.Mkdir(filepath.Dir(p), os.ModePerm)
		}
	}
	return p
}

func LookPath(file string) string {
	path := os.Getenv("path")
	for _, dir := range filepath.SplitList(path) {
		if FileExists(filepath.Join(dir, file)) {
			return filepath.Join(dir, file)
		}
	}
	return ""
}
