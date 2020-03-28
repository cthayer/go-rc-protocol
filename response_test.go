package rc_protocol

import (
	"reflect"
	"testing"
)

func TestNewResponse(t *testing.T) {
	want := response{ExitCode: -1}

	if got := newResponse(""); !reflect.DeepEqual(got, want) {
		t.Errorf("newResponse() = %q, want %q", got, want)
	}
}
