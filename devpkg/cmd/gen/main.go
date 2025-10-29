package main

import (
	"context"
	"os"
	"path/filepath"

	_ "github.com/xoctopus/enumx/devpkg/enumx"
	_ "github.com/xoctopus/genx/devpkg/errorx"
	"github.com/xoctopus/genx/pkg/genx"
	"github.com/xoctopus/x/misc/must"
	_ "github.com/xoctopus/x/reflectx"
	_ "github.com/xoctopus/x/textx"
)

func main() {
	cwd := must.NoErrorV(os.Getwd())

	must.NoError(
		genx.NewContext(&genx.Args{
			Entrypoint: []string{filepath.Join(cwd, "consts")},
		}).Execute(context.Background(), genx.Get()...),
	)
}
