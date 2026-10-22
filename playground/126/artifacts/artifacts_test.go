package artifacts_test

import (
	"os"
	"path/filepath"
	"testing"
)

func TestArtifactDir(t *testing.T) {
	dir := t.ArtifactDir()
	if dir == "" {
		t.Fatal("empty artifact dir")
	}

	path := filepath.Join(dir, "note.txt")
	if err := os.WriteFile(path, []byte("go1.26"), 0o644); err != nil {
		t.Fatal(err)
	}
	b, err := os.ReadFile(path)
	if err != nil || string(b) != "go1.26" {
		t.Fatalf("read back: %q err=%v", b, err)
	}
}
