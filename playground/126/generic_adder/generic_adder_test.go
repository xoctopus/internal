package generic_adder_test

import "testing"

type Adder[A Adder[A]] interface {
	Add(A) A
}

type IntAdder int

func (a IntAdder) Add(b IntAdder) IntAdder { return a + b }

func sum[A Adder[A]](x, y A) A { return x.Add(y) }

func TestGenericSelfReferenceInConstraint(t *testing.T) {
	if got := sum(IntAdder(2), IntAdder(40)); got != 42 {
		t.Fatalf("got %d", got)
	}
}
