package test

import (
	"bytes"
	"path"
	"runtime"
	"strings"

	"github.com/sergi/go-diff/diffmatchpatch"
)

func RootDir() string {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		panic("test failed to get file name")
	}

	rootDir := path.Join(path.Dir(filename), "..")

	return rootDir
}

func GetColoredDiff(a, b string) string {
	dmp := diffmatchpatch.New()

	diffs := dmp.DiffMain(a, b, true)

	colors := map[diffmatchpatch.Operation]string{
		diffmatchpatch.DiffInsert: "\x1b[42m",
		diffmatchpatch.DiffDelete: "\x1b[41m",
	}

	var buff bytes.Buffer
	for _, diff := range diffs {
		text := diff.Text

		if diff.Type == diffmatchpatch.DiffEqual {
			buff.WriteString(text)
			continue
		}

		segments := strings.Split(text, "\n")
		for i, s := range segments {
			if s != "" {
				buff.WriteString(colors[diff.Type])
				buff.WriteString(s)
				buff.WriteString("\x1b[0m")
			}

			if i < len(segments)-1 {
				buff.WriteString(colors[diff.Type])
				buff.WriteString("\\n")
				buff.WriteString("\x1b[0m")

				if diff.Type == diffmatchpatch.DiffInsert {
					buff.WriteString("\n")
				}
			}
		}
	}

	return buff.String()
}
