package main

import (
	_ "github.com/xoctopus/genx/devpkg/docx"

	"github.com/xoctopus/internal/devpkg/devtool"
)

func main() {
	cmd := devtool.Command()
	if err := cmd.Execute(); err != nil {
		cmd.Println(err)
		return
	}
}
