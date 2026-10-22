package jsontext

import (
	jsonv1 "encoding/json"
	jsonv2 "encoding/json/v2"
	"fmt"
	"testing"
)

type benchNested struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type benchPayload struct {
	A int         `json:"a"`
	B string      `json:"b"`
	C []int       `json:"c"`
	D benchNested `json:"d"`
}

type benchRecord struct {
	ID   int          `json:"id"`
	Name string       `json:"name"`
	Tags []string     `json:"tags"`
	Meta benchPayload `json:"meta"`
}

var (
	sampleSmall = []byte(`{
        "a": 123,
        "b": "hello",
        "c": [1,2,3,4,5],
        "d": {"x": 9, "y": 10}
    }`)

	sampleLarge []byte
)

func init() {
	records := make([]benchRecord, 800)
	for i := range records {
		records[i] = benchRecord{
			ID:   i,
			Name: fmt.Sprintf("item-%d", i),
			Tags: []string{"alpha", "beta", "gamma"},
			Meta: benchPayload{
				A: i,
				B: "hello",
				C: []int{1, 2, 3, 4, 5},
				D: benchNested{X: 9, Y: 10},
			},
		}
	}
	var err error
	sampleLarge, err = jsonv1.Marshal(records)
	if err != nil {
		panic(err)
	}
}

func BenchmarkUnmarshal(b *testing.B) {
	b.Run("Small/MapAny", func(b *testing.B) {
		benchmarkUnmarshalMap(b, sampleSmall)
	})
	b.Run("Small/Struct", func(b *testing.B) {
		benchmarkUnmarshalStruct(b, sampleSmall, func() any { return new(benchPayload) })
	})
	b.Run("Large/SliceAny", func(b *testing.B) {
		benchmarkUnmarshalSliceAny(b, sampleLarge)
	})
	// Reuse: 反复 Unmarshal 进同一个 *[]benchRecord. V2 可能为扩容多分配, 与 Fresh 对比着看.
	b.Run("Large/Struct/Reuse", func(b *testing.B) {
		benchmarkUnmarshalStruct(b, sampleLarge, func() any { return new([]benchRecord) })
	})
	b.Run("Large/Struct/Fresh", func(b *testing.B) {
		b.Run("V1", func(b *testing.B) {
			for b.Loop() {
				var v []benchRecord
				if err := jsonv1.Unmarshal(sampleLarge, &v); err != nil {
					b.Fatal(err)
				}
			}
		})
		b.Run("V2", func(b *testing.B) {
			for b.Loop() {
				var v []benchRecord
				if err := jsonv2.Unmarshal(sampleLarge, &v); err != nil {
					b.Fatal(err)
				}
			}
		})
	})
}

func benchmarkUnmarshalSliceAny(b *testing.B, data []byte) {
	b.Run("V1", func(b *testing.B) {
		var v []any
		for b.Loop() {
			if err := jsonv1.Unmarshal(data, &v); err != nil {
				b.Fatal(err)
			}
		}
	})
	b.Run("V2", func(b *testing.B) {
		var v []any
		for b.Loop() {
			if err := jsonv2.Unmarshal(data, &v); err != nil {
				b.Fatal(err)
			}
		}
	})
}

func benchmarkUnmarshalMap(b *testing.B, data []byte) {
	b.Run("V1", func(b *testing.B) {
		var v map[string]any
		for b.Loop() {
			if err := jsonv1.Unmarshal(data, &v); err != nil {
				b.Fatal(err)
			}
		}
	})
	b.Run("V2", func(b *testing.B) {
		var v map[string]any
		for b.Loop() {
			if err := jsonv2.Unmarshal(data, &v); err != nil {
				b.Fatal(err)
			}
		}
	})
}

func benchmarkUnmarshalStruct(b *testing.B, data []byte, newDest func() any) {
	b.Run("V1", func(b *testing.B) {
		dest := newDest()
		for b.Loop() {
			if err := jsonv1.Unmarshal(data, dest); err != nil {
				b.Fatal(err)
			}
		}
	})
	b.Run("V2", func(b *testing.B) {
		dest := newDest()
		for b.Loop() {
			if err := jsonv2.Unmarshal(data, dest); err != nil {
				b.Fatal(err)
			}
		}
	})
}

func BenchmarkMarshal(b *testing.B) {
	mapVal := map[string]any{
		"a": 123,
		"b": "hello",
		"c": []int{1, 2, 3, 4, 5},
		"d": map[string]any{"x": 9, "y": 10},
	}
	structVal := benchPayload{
		A: 123,
		B: "hello",
		C: []int{1, 2, 3, 4, 5},
		D: benchNested{X: 9, Y: 10},
	}
	largeStruct := make([]benchRecord, 800)
	for i := range largeStruct {
		largeStruct[i] = benchRecord{
			ID:   i,
			Name: fmt.Sprintf("item-%d", i),
			Tags: []string{"alpha", "beta", "gamma"},
			Meta: structVal,
		}
	}

	b.Run("Small/MapAny", func(b *testing.B) {
		benchmarkMarshal(b, mapVal)
	})
	b.Run("Small/Struct", func(b *testing.B) {
		benchmarkMarshal(b, structVal)
	})
	b.Run("Large/Struct", func(b *testing.B) {
		benchmarkMarshal(b, largeStruct)
	})
}

func benchmarkMarshal(b *testing.B, v any) {
	b.Run("V1", func(b *testing.B) {
		for b.Loop() {
			if _, err := jsonv1.Marshal(v); err != nil {
				b.Fatal(err)
			}
		}
	})
	b.Run("V2", func(b *testing.B) {
		for b.Loop() {
			if _, err := jsonv2.Marshal(v); err != nil {
				b.Fatal(err)
			}
		}
	})
}
