package main

import (
	_ "github.com/xoctopus/genx/devpkg/docx"

	"github.com/xoctopus/internal/devpkg/devtool"
)

func init() {
	// cwd := must.NoErrorV(os.Getwd())

	// ctx := genx.NewContext(&genx.Args{Entrypoint: []string{filepath.Join(cwd, "devpkg", "...")}})

	// if err := ctx.Execute(context.Background(), genx.Get()...); err != nil {
	// 	panic(err)
	// }
}

func main() {
	cmd := devtool.Command()
	if err := cmd.Execute(); err != nil {
		cmd.Println(err)
		return
	}
}
