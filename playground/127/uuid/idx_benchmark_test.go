package uuid_test

import (
	"testing"
	stduuid "uuid"

	guuid "github.com/google/uuid"
	"github.com/oklog/ulid/v2"
	"github.com/xoctopus/sfid/pkg/sfid"
)

func BenchmarkIDGen(b *testing.B) {
	g := sfid.NewDefaultIDGen(1)

	b.Run("xoctopus/sfid", func(b *testing.B) {
		for b.Loop() {
			_, _ = g.ID()
		}
	})

	b.Run("uuid/stdlib/v4", func(b *testing.B) {
		for b.Loop() {
			_ = stduuid.NewV4()
		}
	})

	b.Run("uuid/stdlib/v7", func(b *testing.B) {
		for b.Loop() {
			_ = stduuid.NewV7()
		}
	})

	b.Run("uuid/google/v4", func(b *testing.B) {
		for b.Loop() {
			_ = guuid.New()
		}
	})

	b.Run("uuid/google/v7", func(b *testing.B) {
		for b.Loop() {
			_, _ = guuid.NewV7()
		}
	})

	b.Run("ulid/v2", func(b *testing.B) {
		for b.Loop() {
			_ = ulid.Make()
		}
	})
}
