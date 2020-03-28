package rc_protocol

import (
	"crypto"
	"crypto/sha256"
	"encoding/base64"
	"reflect"
	"regexp"
	"testing"
)

func TestNewAuth(t *testing.T) {
	want := auth{
		headerPattern: regexp.MustCompile(`^RC [^\s;]+;[^\s;]+;.+$`),
		sigAlgorithm: sha256.New(),
		sigAlgorithmCrypto: crypto.SHA256,
		sigEncoding: base64.StdEncoding,
	}

	if got := newAuth(); !reflect.DeepEqual(got, want) {
		t.Errorf("newAuth() = %q, want %v", got, want)
	}
}
