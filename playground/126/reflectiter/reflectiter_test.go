package reflectiter_test

import (
	"reflect"
	"testing"
)

type Point struct {
	X, Y int
}

func TestTypeFieldsIterator(t *testing.T) {
	rt := reflect.TypeFor[Point]()

	var names []string
	for f := range rt.Fields() {
		names = append(names, f.Name)
	}
	if len(names) != 2 || names[0] != "X" || names[1] != "Y" {
		t.Fatalf("Fields: %v", names)
	}

	rf := reflect.ValueOf(Point{3, 4})
	var vals []int
	for f, v := range rf.Fields() {
		if f.Name == "X" || f.Name == "Y" {
			vals = append(vals, int(v.Int()))
		}
	}
	if len(vals) != 2 || vals[0] != 3 || vals[1] != 4 {
		t.Fatalf("Value.Fields: %v", vals)
	}
}
