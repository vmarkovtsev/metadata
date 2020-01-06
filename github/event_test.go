package github

import (
	"bytes"
	"testing"

	fuzz "github.com/google/gofuzz"
	"github.com/stretchr/testify/require"
)

func TestEvent(t *testing.T) {
	require := require.New(t)

	var src Event
	fuzz.New().Fuzz(&src)

	data, err := MarshalEvent(&src)
	require.NoError(err)

	dst, err := UnmarshalEvent(data)
	require.NoErrorf(err, "event: %v", src)

	require.Equal(src.Type, (*dst).Type)
	require.Truef(bytes.Equal(src.Payload, (*dst).Payload), "event: %v, expected: %v", *dst, src)
}
