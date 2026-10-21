package weakptr

import (
	"fmt"
	"weak"
)

// EventHub 事件中心, 按发布-订阅模型广播事件.
// 订阅者使用 weak.Pointer[Listener] 引用, 而非 *Subscriber.
//
// 使用场景: 主题需要记住「谁要收通知」, 但 Subscriber 的生命周期应由调用方决定,
// 不应因为 Register 过一次就被 EventHub 一直拽在内存里.
// weak.Make 只登记弱引用, 不阻止 GC; Notify 时对 Value() 为 nil 的条目跳过并移出列表.
//
// 与 []weak.Pointer 相对的是 []*Subscriber: 强引用下若从不 Unregister 会泄漏订阅者.
// 标准库也注明 weak 更偏 runtime/基础设施; 业务里若要求可预期的注销, 显式 Unregister 通常更清晰.
type EventHub struct {
	observers []weak.Pointer[Subscriber]
}

func (e *EventHub) Register(o *Subscriber) {
	e.observers = append(e.observers, weak.Make(o))
	fmt.Printf("%s registered\n", o.ID())
}

func (e *EventHub) Notify(content string) {
	alive := e.observers[:0]
	for _, wp := range e.observers {
		o := wp.Value()
		if o != nil {
			o.Notify(content)
			alive = append(alive, wp)
		}
	}
	e.observers = alive
}

func NewSubscriber(id string) *Subscriber {
	return &Subscriber{id: id}
}

type Subscriber struct {
	id string
}

func (l *Subscriber) ID() string {
	return l.id
}

func (l *Subscriber) Notify(content string) {
	fmt.Printf("%s received: %s\n", l.id, content)
}
