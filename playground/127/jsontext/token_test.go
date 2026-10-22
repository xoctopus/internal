package jsontext

import (
	"bytes"
	"encoding/json/jsontext"
	"testing"
)

func TestJsonText(t *testing.T) {
	dec := jsontext.NewDecoder(bytes.NewReader(sampleSmall))
	token, _ := dec.ReadToken()
	t.Log(token.String())
	v, _ := dec.ReadValue()
	t.Log(v.String())

	token, _ = dec.ReadToken()
	t.Log(token.String())
}
