package devtool

import (
	"github.com/spf13/cobra"
	"github.com/xoctopus/confx/pkg/cmdx"

	"github.com/xoctopus/internal/devpkg/devtool/internal"
)

var root = &cobra.Command{}

func init() {
	root.AddCommand(cmdx.NewCommand("make", &internal.Makefile{}).Cmd())
}

func Command() *cobra.Command {
	return root
}
