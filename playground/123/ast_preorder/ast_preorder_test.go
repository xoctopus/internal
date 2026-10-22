package ast_preorder_test

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"strings"

	"golang.org/x/tools/go/packages"
)

const path = "github.com/xoctopus/internal/playground/123/ast_preorder"

var (
	dataStruct *ast.StructType
	typesInfo  *types.Info
)

func init() {
	pkgs, err := packages.Load(&packages.Config{
		Mode: packages.NeedSyntax | packages.NeedTypes | packages.NeedTypesInfo,
	}, path)
	if err != nil {
		panic(err)
	}
	if len(pkgs) != 1 {
		panic("expected one package")
	}
	pkg := pkgs[0]
	if len(pkg.Errors) != 0 {
		panic(pkg.Errors[0])
	}
	dataStruct = structType(pkg, "Data")
	typesInfo = pkg.TypesInfo
}

func structType(pkg *packages.Package, name string) *ast.StructType {
	for _, f := range pkg.Syntax {
		for _, decl := range f.Decls {
			gd, ok := decl.(*ast.GenDecl)
			if !ok || gd.Tok != token.TYPE {
				continue
			}
			for _, spec := range gd.Specs {
				ts, ok := spec.(*ast.TypeSpec)
				if !ok || ts.Name.Name != name {
					continue
				}
				st, ok := ts.Type.(*ast.StructType)
				if !ok {
					panic("Data is not a struct")
				}
				return st
			}
		}
	}
	panic("type " + name + " not found")
}

func ExamplePreorder() {
	for n := range ast.Preorder(dataStruct) {
		if f, ok := n.(*ast.Field); ok {
			names := make([]string, 0, len(f.Names))
			for _, ident := range f.Names {
				names = append(names, ident.Name)
			}
			name := strings.Join(names, ", ")
			if name == "Z" {
				break
			}
			typ := typesInfo.Types[f.Type]

			fmt.Printf("Field: %s %s\n", name, typ.Type.String())
		}
	}

	ast.Inspect(dataStruct, func(n ast.Node) bool {
		if f, ok := n.(*ast.Field); ok {
			names := make([]string, 0, len(f.Names))
			for _, ident := range f.Names {
				names = append(names, ident.Name)
			}
			name := strings.Join(names, ", ")
			if name == "Z" {
				return false
			}
			typ := typesInfo.Types[f.Type]

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
