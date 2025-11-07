package week_ptr

import (
	"fmt"
	"runtime"
	"testing"
	"weak"
)

// weak.Make if you need to cache an object without affect its lifetime, use it.
// NOTE: weak.Pointer are for system builders, not for general application code.
// eg:
//
//	make[string]*Object              // user control Object's lifecycle
//	make[string]weak.Pointer[Object] // go gc control Object's lifecycle
func TestWeekPtr(t *testing.T) {
	v := new(int)
	*v = 100
	wp := weak.Make(v)

	t.Log(wp.Value() == nil) // expect false

	v = nil
	t.Log(wp.Value() == nil) // expect false

	runtime.GC()
	if wp.Value() == nil {
		t.Log("v has be gc")
	}
}

type Events struct {
	observers []weak.Pointer[Listener]
}

func (e *Events) Register(o *Listener) {
	e.observers = append(e.observers, weak.Make(o))
}

func (e *Events) Notify(content string) {
	alive := e.observers[:0]
	for _, wp := range e.observers {
		o := wp.Value()
		if o != nil {
			o.OnEvent(content)
			alive = append(alive, wp)
		}
	}
	e.observers = alive
}

type Observer[T any] interface {
	ID() string
	OnEvent(string)
	Underlying() weak.Pointer[T]
}

type Listener struct {
	id string
}

func (l *Listener) ID() string {
	return l.id
}

func (l *Listener) OnEvent(content string) {
	fmt.Printf("%s received: %s\n", l.id, content)
}

// Observer interface is not used. So think how to avoid an interface to be garbage collected
func Example_weak_pointer() {
	events := &Events{}

	l1 := &Listener{id: "1"}
	l2 := &Listener{id: "2"}
	events.Register(l1)
	events.Register(l2)

	events.Notify("some before gc")

	l1 = nil
	runtime.GC()
	runtime.KeepAlive(l2)

	events.Notify("some after gc")

	// Output:
	// 1 received: some before gc
	// 2 received: some before gc
	// 2 received: some after gc
}
