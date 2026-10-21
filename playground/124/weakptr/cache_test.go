package weakptr_test

import (
	"fmt"
	"runtime"

	"github.com/xoctopus/internal/playground/124/weakptr"
)

func ExampleCache_weakEntryExpiresAfterGC() {
	c := &weakptr.Cache{}

	v := &weakptr.Big{Label: "parsed-doc"}
	c.Put("doc", v)

	if got := c.Get("doc"); got != nil {
		fmt.Println("hit:", got.Label)
	}

	// 调用方不再持有 *Big; 缓存里只剩 weak, 不阻止 GC.
	// 在内存吃紧的情况下 weak 并不因为 Cache 的引用阻止 GC.
	v = nil
	runtime.GC()

	if c.Get("doc") == nil {
		fmt.Println("miss after gc")
	}

	// Output:
	// hit: parsed-doc
	// miss after gc
}
