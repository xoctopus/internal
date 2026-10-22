package hpke_test

import (
	"crypto/ecdh"
	"crypto/hpke"
	"crypto/rand"
	"testing"
)

func TestSealOpen(t *testing.T) {
	curve := ecdh.P256()
	priv, err := curve.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}

	pk, err := hpke.NewDHKEMPublicKey(priv.PublicKey())
	if err != nil {
		t.Fatal(err)
	}
	sk, err := hpke.NewDHKEMPrivateKey(priv)
	if err != nil {
		t.Fatal(err)
	}

	kdf := hpke.HKDFSHA256()
	aead := hpke.AES128GCM()
	info := []byte("playground/126")

	plain := []byte("hello hpke")
	ct, err := hpke.Seal(pk, kdf, aead, info, plain)
	if err != nil {
		t.Fatal(err)
	}

	got, err := hpke.Open(sk, kdf, aead, info, ct)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != string(plain) {
		t.Fatalf("got %q", got)
	}
}
