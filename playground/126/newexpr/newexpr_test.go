package newexpr_test

import (
	"encoding/json"
	"testing"
	"time"
)

type Person struct {
	Name string `json:"name"`
	Age  *int   `json:"age"`
}

func yearsSince(t time.Time) int {
	return int(time.Since(t).Hours() / (365.25 * 24))
}

func TestNewExpressionOperand(t *testing.T) {
	v := new(int64(300))
	if *v != 300 {
		t.Fatalf("got %d", *v)
	}

	born := time.Now().Add(-10 * 365 * 24 * time.Hour)
	data, err := json.Marshal(Person{
		Name: "Ada",
		Age:  new(yearsSince(born)),
	})
	if err != nil {
		t.Fatal(err)
	}
	if string(data) == `{"name":"Ada","age":null}` {
		t.Fatal("expected age field populated")
	}
}
