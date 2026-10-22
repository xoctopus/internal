package errorsastype_test

import (
	"errors"
	"fmt"
	"testing"
)

type codedError struct {
	Code int
}

func (e codedError) Error() string { return fmt.Sprintf("code=%d", e.Code) }

func TestErrorsAsType(t *testing.T) {
	err := fmt.Errorf("wrap: %w", codedError{Code: 126})

	ce, ok := errors.AsType[codedError](err)
	if !ok {
		t.Fatal("AsType failed")
	}
	if ce.Code != 126 {
		t.Fatalf("got code %d", ce.Code)
	}

	_, ok = errors.AsType[codedError](errors.New("other"))
	if ok {
		t.Fatal("expected false for unrelated error")
	}
}
