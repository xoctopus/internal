package internal

import (
	"io"
	"os"
	"path/filepath"
	"strings"
)

func IsGoProject(wd string) bool {
	tar := filepath.Join(wd, "go.mod")
	i, err := os.Stat(tar)
	return err == nil && !i.IsDir()
}

func IsGitRepository(wd string) bool {
	tar := filepath.Join(wd, ".git")
	i, err := os.Stat(tar)
	return err == nil && i.IsDir()
}

func WriteKeyValAlign(w io.Writer, prefix, join string, kvs [][2]string) error {
	maxlen := 0
	for _, kv := range kvs {
		if len(kv[0]) > maxlen {
			maxlen = len(kv[0])
		}
	}
	for _, kv := range kvs {
		k := prefix + kv[0] + strings.Repeat(" ", maxlen-len(kv[0]))
		v := kv[1]
		if _, err := w.Write([]byte(k)); err != nil {
			return err
		}
		if _, err := w.Write([]byte(" " + join + " ")); err != nil {
			return err
		}
		if _, err := w.Write([]byte(v + "\n")); err != nil {
			return err
		}
	}
	return nil
}

var (
	DefaultTestIgnores = []string{
		"_gen.go",
		".pb.go",
		"_mock.go",
		"_genx_",
		"main.go",
		"testing.go",
		"example/",
	}
	DefaultFormatIgnores = []string{
		".git/",
		".xgo/",
		"*.pb.go",
		"*_genx_*",
	}
	DefaultGoTools = [][2]string{
		{"GOTEST", "go"},
		{"GOBUILD", "go"},
	}
)
