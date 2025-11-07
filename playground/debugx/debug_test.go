package debugx

import (
	"fmt"
	"os"
	"path/filepath"
)

func _ExampleDebugInfo() {
	cwd, _ := os.Getwd()
	filename := filepath.Join(cwd, "testdata", "main")

	fmt.Println(filename)
	DebugInfo(filename)

	// Output:
}
