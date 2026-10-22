package slogmulti_test

import (
	"bytes"
	"log/slog"
	"strings"
	"testing"
)

func TestNewMultiHandler(t *testing.T) {
	var a, b bytes.Buffer
	h := slog.NewMultiHandler(
		slog.NewTextHandler(&a, nil),
		slog.NewJSONHandler(&b, nil),
	)
	l := slog.New(h)
	l.Info("hello", "ver", 126)

	if !strings.Contains(a.String(), "hello") {
		t.Fatalf("text handler: %q", a.String())
	}
	if !strings.Contains(b.String(), `"ver":126`) {
		t.Fatalf("json handler: %q", b.String())
	}
}
