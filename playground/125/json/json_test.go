package json_test

import (
	"encoding/json"
	"testing"
)

func TestJsonUnexported(t *testing.T) {
	type embed struct{ X int }
	type S struct{ *embed }

	i := &S{&embed{5}}
	o := &S{}

	b, _ := json.Marshal(i)
	t.Log(string(b))

	err := json.Unmarshal(b, o)
	t.Log(err)
	// if err == nil {
	// 	defer func() {
	// 		t.Log(recover())
	// 	}()
	// 	t.Log(o.X)
	// }
}
