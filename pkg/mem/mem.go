package mem

import (
	"fmt"
	"runtime"
)

func OutputUsage(label string) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	_, _ = fmt.Printf("%s:\n", label)
	_, _ = fmt.Printf("\tAlloc      = %v MiB\n", bToMb(m.Alloc))
	_, _ = fmt.Printf("\tTotalAlloc = %v MiB\n", bToMb(m.TotalAlloc))
	_, _ = fmt.Printf("\tSys        = %v MiB\n", bToMb(m.Sys))
	_, _ = fmt.Printf("\tHeapAlloc  = %v MiB\n", bToMb(m.HeapAlloc))
	_, _ = fmt.Printf("\tHeapSys    = %v MiB\n", bToMb(m.HeapSys))
	_, _ = fmt.Printf("\tHeapInuse  = %v MiB\n", bToMb(m.HeapInuse))
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
