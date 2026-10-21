package weakptr_test

import (
	"runtime"

	"github.com/xoctopus/internal/playground/124/weakptr"
)

func ExampleEventHub_notifyAfterSubscriberGC() {
	hub := &weakptr.EventHub{}

	sub1 := weakptr.NewSubscriber("sub1")
	sub2 := weakptr.NewSubscriber("sub2")
	hub.Register(sub1)
	hub.Register(sub2)

	// 执行通知
	hub.Notify("before gc")

	// 此时 sub2 生命周期已结束 (资源被GC回收)
	sub2 = nil
	runtime.GC()
	runtime.KeepAlive(sub1)

	// 执行通知
	hub.Notify("after gc")

	// Output:
	// sub1 registered
	// sub2 registered
	// sub1 received: before gc
	// sub2 received: before gc
	// sub1 received: after gc
}
