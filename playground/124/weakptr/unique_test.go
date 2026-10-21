package weakptr_test

import (
	"fmt"
	"runtime"

	"github.com/xoctopus/internal/playground/124/weakptr"
)

func ExampleInternTable_reinternAfterGC() {
	tab := &weakptr.InternTable{}

	a := tab.Intern("east")
	b := tab.Intern("east")
	fmt.Println("same instance:", a == b)
	fmt.Println("gen while alive:", a.Gen)

	// 没有任何 *Symbol 强引用时, table 里的 weak 拦不住 GC.
	a, b = nil, nil
	runtime.GC()

	c := tab.Intern("east")
	fmt.Println("gen after gc:", c.Gen)

	// Output:
	// same instance: true
	// gen while alive: 1
	// gen after gc: 2
}
