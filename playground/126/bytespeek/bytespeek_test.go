package bytespeek_test

import (
	"bytes"
	"testing"
)

func TestBufferPeek(t *testing.T) {
	var buf bytes.Buffer
	_, _ = buf.WriteString("hello")

	head, err := buf.Peek(3)
	if err != nil {
		t.Fatal(err)
	}
	if string(head) != "hel" {
		t.Fatalf("Peek: got %q", head)
	}

	all, err := buf.Peek(5)
	if err != nil {
		t.Fatal(err)
	}
	if string(all) != "hello" {
		t.Fatalf("Peek all: got %q", all)
	}

	_, _ = buf.ReadByte()
	if rest, _ := buf.Peek(10); string(rest) != "ello" {
		t.Fatalf("after ReadByte Peek: got %q", rest)
	}
}
