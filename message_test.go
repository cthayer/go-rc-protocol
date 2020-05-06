package rc_protocol

import (
	"reflect"
	"testing"
)

func TestNewMessage(t *testing.T) {
	want := Message{}

	if got := newMessage(""); !reflect.DeepEqual(got, want) {
		t.Errorf("newMessage() = %q, want %q", got, want)
	}
}
