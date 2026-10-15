package debugx

import (
	"os"
	"path/filepath"
)

func ExampleDebugInfo() {
	cwd, _ := os.Getwd()
	filename := filepath.Join(cwd, "testdata", "main")

	// fmt.Println(filename)
	DebugInfo(filename)

	// Output:
}
