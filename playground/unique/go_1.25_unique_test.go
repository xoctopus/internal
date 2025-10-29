package unique_test

import (
	"testing"

	"github.com/xoctopus/internal/playground/unique"
)

func TestUnique(t *testing.T) {
	unique.Unique()
}

func TestWithoutUnique(t *testing.T) {
	unique.WithoutUnique()
}
