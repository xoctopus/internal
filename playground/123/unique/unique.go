package unique

import (
	"runtime"
	"strings"
	"unique"
)

const (
	SCALE          = 1000000
	STRING_LENGTH  = 20
	DISTINCT_COUNT = 1000
)

var seed = make([]string, DISTINCT_COUNT)

func init() {
	for i := range DISTINCT_COUNT {
		seed[i] = strings.Repeat(string(rune('A'+i%26)), STRING_LENGTH)
	}
}

func Unique() {
	data := make([]unique.Handle[string], SCALE)
	for i := range SCALE {
		data[i] = unique.Make(seed[i%DISTINCT_COUNT])
	}

	runtime.GC()

	OutputUsage("Unique")

	runtime.KeepAlive(data)
}

func WithoutUnique() {
	data := make([]string, SCALE)
	for i := range SCALE {
		data[i] = data[i%DISTINCT_COUNT]
	}

	runtime.GC()

	OutputUsage("WithoutUnique")

	runtime.KeepAlive(data)
}
