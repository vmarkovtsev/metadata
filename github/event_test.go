package github

import (
	"bytes"
	"testing"

	fuzz "github.com/google/gofuzz"
)

func TestEvent(t *testing.T) {
	var src Event
	fuzz.New().Fuzz(&src)

	data, err := MarshalEvent(&src)
	if err != nil {
		t.Fatal(err)
	}

	dst, err := UnmarshalEvent(data)
	if err != nil {
		t.Errorf("event: %v,  error: %s", src, err)
	}

	if src.Type != (*dst).Type || !bytes.Equal(src.Payload, (*dst).Payload) {
		t.Errorf("event: %v, expected: %v", *dst, src)
	}
}
