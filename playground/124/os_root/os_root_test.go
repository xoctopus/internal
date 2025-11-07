package os_root

import (
	"os"
	"testing"
)

// os.Root presents a rooted directory context. Without changing the host
// permissions, provide Go programs with a safe, controlled, non-escapable
// directory scope that confines all filesystem operations to a specific
// directory tree.
func TestOsRoot(t *testing.T) {
	root, err := os.OpenRoot("/tmp")
	if err != nil {
		t.Fatal(err)
	}
	defer root.Close()

	f, err := root.Create("temp.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	err = root.MkdirAll("folder", 0o777)
	if err != nil {
		t.Fatal(err)
	}

	f2, err := root.Create("folder/tmp.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer f2.Close()

	_, err = root.Stat("temp.txt")
	if err != nil {
		t.Fatal(err)
	}

	_, err = root.Stat("folder/tmp.txt")
	t.Log(err)

	_ = root.Remove("temp.txt")

	_, err = root.Stat("/folder/tmp.txt")
	t.Log(err)

	_, err = root.Stat("/etc/passwd")
	t.Log(err)
}
