//go:build jsonv2

package jsontext_test

import (
	jsonv1 "encoding/json"
	jsonv2 "encoding/json/v2"
	"testing"
)

func TestJson(t *testing.T) {
	bad := []byte(`{"a": 1, "b":}`)

	var o1 map[string]any
	t.Log(jsonv1.Unmarshal(bad, &o1))

	// New
	var o2 map[string]any
	t.Log(jsonv2.Unmarshal(bad, &o2))
}
