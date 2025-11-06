package ast_preorder_test

import (
	"context"
	"fmt"
	"go/ast"
	"strings"

	"github.com/xoctopus/pkgx"
)

var (
	pkg  pkgx.Package
	path = "github.com/xoctopus/internal/playground/123/ast_preorder"
)

func init() {
	pkgs := pkgx.NewPackages(context.Background(), path)
	pkg = pkgs.Package(path)
}

func ExamplePreorder() {
	tt := pkg.TypeNames().ElementByName("Data")

	for n := range ast.Preorder(tt.Node()) {
		if f, ok := n.(*ast.Field); ok {
			names := make([]string, 0, len(f.Names))
			for _, ident := range f.Names {
				names = append(names, ident.Name)
			}
			name := strings.Join(names, ", ")
			if name == "Z" {
				break
			}
			typ := pkg.GoPackage().TypesInfo.Types[f.Type]

			fmt.Printf("Field: %s %s\n", name, typ.Type.String())
		}
	}

	ast.Inspect(tt.Node(), func(n ast.Node) bool {
		if f, ok := n.(*ast.Field); ok {
			names := make([]string, 0, len(f.Names))
			for _, ident := range f.Names {
				names = append(names, ident.Name)
			}
			name := strings.Join(names, ", ")
			if name == "Z" {
				return false
			}
			typ := pkg.GoPackage().TypesInfo.Types[f.Type]

			fmt.Printf("Field: %s %s\n", name, typ.Type.String())
		}
		return true
	})

	// Output:
	// Field: V uint32
	// Field: X uint16
	// Field: Y uint16
	// Field: V uint32
	// Field: X uint16
	// Field: Y uint16
}
