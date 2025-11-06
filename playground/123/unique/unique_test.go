package unique_test

import (
	"testing"

	"github.com/xoctopus/internal/playground/123/unique"
)

func TestUnique(t *testing.T) {
	unique.Unique()
}

func TestWithoutUnique(t *testing.T) {
	unique.WithoutUnique()
}
