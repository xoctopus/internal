// Package jsontext exercises encoding/json/jsontext and json/v2.
//
// Benchmarks (V1 vs V2, small/large, map vs struct):
//
//	go test -tags jsonv2 -bench=. -benchmem ./...
//	go test -tags jsonv2 -bench=. -benchmem -count=6 ./... | benchstat
package jsontext
