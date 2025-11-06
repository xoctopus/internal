package range_over_func

import (
	"fmt"
	"strconv"
	"sync"
	"testing"
)

func RangeInteger() {
	for i := range 5 {
		fmt.Println(i)
	}
}

func RangeString() {
	s := "abc"
	for i := range s {
		fmt.Println(s[i : i+1])
	}
}

func RangeOverFunc() {
	s := sync.Map{}
	for i := range 5 {
		s.Store(i, strconv.Itoa(i))
	}

	for k, v := range s.Range {
		fmt.Printf("%d:%s\n", k, v)
	}
}

func TestRange(t *testing.T) {
	RangeInteger()
	RangeString()
	RangeOverFunc()
}
