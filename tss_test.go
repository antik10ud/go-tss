package tss

import (
	"testing"
	//"crypto/rand"
	"bytes"
	"encoding/hex"
)

func TestKat(t *testing.T) {
	var secret, _ = hex.DecodeString("7465737400")
	var shares = make([][]byte, 2)
	shares[0], _ = hex.DecodeString("01B9FA07E185")
	shares[1], _ = hex.DecodeString("02F5409B4511")
	var recoveredSecret, err = RecoverSecret(shares)
	if err != nil {
		t.Error(err)
	}
	if bytes.Compare(recoveredSecret, secret) != 0 {
		t.Errorf("secret mismatch %x, want %x", recoveredSecret, secret)
	}
}
