package jsontext

import (
	"bytes"
	jsonv1 "encoding/json"
	"encoding/json/jsontext"
	jsonv2 "encoding/json/v2"
	"testing"
)

var sample = []byte(`{
        "a": 123,
        "b": "hello",
        "c": [1,2,3,4,5],
        "d": {"x": 9, "y": 10}
    }`)

func Benchmark_Unmarshal(b *testing.B) {
	var v map[string]any

	b.Run("V1", func(b *testing.B) {
		for b.Loop() {
			_ = jsonv1.Unmarshal(sample, &v)
		}
	})
	b.Run("V2", func(b *testing.B) {
		for b.Loop() {
			_ = jsonv2.Unmarshal(sample, &v)
		}
	})
}

func Benchmark_Marshal(b *testing.B) {
	v := map[string]any{
		"a": 123,
		"b": "hello",
		"c": []int{1, 2, 3, 4, 5},
		"d": map[string]any{"x": 9, "y": 10},
	}

	b.Run("V1", func(b *testing.B) {
		_, _ = jsonv1.Marshal(v)
	})
	b.Run("V2", func(b *testing.B) {
		_, _ = jsonv2.Marshal(v)
	})
}

func TestJsonText(t *testing.T) {
	dec := jsontext.NewDecoder(bytes.NewReader(sample))
	token, _ := dec.ReadToken()
	t.Log(token.String())
	v, _ := dec.ReadValue()
	t.Log(v.String())

	token, _ = dec.ReadToken()
	t.Log(token.String())
}
