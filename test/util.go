package test

import (
	"path"
	"runtime"
)

func RootDir() string {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		panic("test failed to get file name")
	}

	rootDir := path.Join(path.Dir(filename), "..")

	return rootDir
}
