package rc_protocol

import (
	"reflect"
	"regexp"
	"testing"
)

func TestNewRCProtocol(t *testing.T) {
	want := protocol{newAuth()}

	if got := NewRCProtocol(); !reflect.DeepEqual(got, want) {
		t.Errorf("NewRCProtocol() = %q, want %q", got, want)
	}
}

func TestProtocol_Message(t *testing.T) {
	want := newMessage("")
	proto := NewRCProtocol()

	if got := proto.Message(""); !reflect.DeepEqual(got, want) {
		t.Errorf("Protocol.Message() = %q, want %q", got, want)
	}
}

func TestProtocol_Response(t *testing.T) {
	want := newResponse("")
	proto := NewRCProtocol()

	if got := proto.Response(""); !reflect.DeepEqual(got, want) {
		t.Errorf("Protocol.Response() = %q, want %q", got, want)
	}
}

func TestAuth_GetHeaderName(t *testing.T) {
	proto := NewRCProtocol()
	want := AUTH_HEADER_NAME

	if got := proto.GetHeaderName(); got != want {
		t.Errorf("Expecting: %q, Got: %q", want, got)
	}
}

func TestProtocol_ParseHeader(t *testing.T) {
	proto := NewRCProtocol()
	want := []string{"user", "2020-01-01T04:01:59-0700", "asdfmbhfefmd8347336;eirey"}
	header := "RC user;2020-01-01T04:01:59-0700;asdfmbhfefmd8347336;eirey"

	if got := proto.ParseHeader(header); !reflect.DeepEqual(got, want) {
		t.Errorf("ParseHeader() = %q, wanted %q", got, want)
	}
}

func TestProtocol_CreateSig(t *testing.T) {
	proto := NewRCProtocol()
	want := regexp.MustCompile(`^RC user;\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}-\d{4};.{50,}={1,2}$`)

	if got, err := proto.CreateSig("user", "test/keys"); !want.MatchString(got) || err != nil {
		t.Errorf("CreateSig() = %q, wanted %q", got, want.String())
		t.Errorf("CreateSig() Error = %q, wanted %v", err, nil)
	}
}

func TestAuth_CheckSig(t *testing.T) {
	proto := NewRCProtocol()

	sig, err := proto.CreateSig("user", "test/keys")

	if err != nil {
		t.Errorf("Failed to create signature to check: %q", err)
		return
	}

	if validSig, err := proto.CheckSig(sig, "test/certs"); !validSig || err != nil {
		t.Errorf("CheckSig() = %v, wanted %v", validSig, true)
		t.Errorf("CheckSig() Error = %q, wanted %v", err, nil)
	}
}
