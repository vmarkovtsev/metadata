package github

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"errors"
	"fmt"
	"hash"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	fuzz "github.com/google/gofuzz"
	"github.com/stretchr/testify/require"
)

const (
	// sha1Prefix is the prefix used by GitHub before the HMAC hexdigest.
	sha1Prefix = "sha1"
	// sha256Prefix and sha512Prefix are provided for future compatibility.
	sha256Prefix = "sha256"
	sha512Prefix = "sha512"
	// signatureHeader is the GitHub header key used to pass the HMAC hexdigest.
	signatureHeader = "X-Hub-Signature"
	// eventTypeHeader is the GitHub header key used to pass the event type.
	eventTypeHeader = "X-Github-Event"
	// deliveryIDHeader is the GitHub header key used to pass the unique ID for the webhook event.
	deliveryIDHeader = "X-Github-Delivery"
)

func TestWebhook(t *testing.T) {
	require := require.New(t)

	tests := []struct {
		name        string
		method      string
		header      map[string]string
		payload     []byte
		reponseCode int
	}{
		{
			name:   "Valid Request",
			method: "POST",
			header: map[string]string{
				"Content-Type":  "application/json",
				eventTypeHeader: "TEST",
			},
			reponseCode: http.StatusOK,
		},
		{
			name:        "Invalid method and missed header",
			method:      "GET",
			header:      map[string]string{},
			reponseCode: http.StatusBadRequest,
		},
		{
			name:        "Missed Content-Type header",
			method:      "POST",
			header:      map[string]string{eventTypeHeader: "TEST"},
			reponseCode: http.StatusBadRequest,
		},
		{
			name:   "Missed X-Github-Event header",
			method: "POST",
			header: map[string]string{
				"Content-Type": "application/json",
			},
			reponseCode: http.StatusBadRequest,
		},
	}

	fuzz := fuzz.New()
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ch := make(chan []byte, 1)
			defer close(ch)

			wh := &Webhook{
				SecretKey: []byte("SECRET_TOKEN"),
				OnEvent: func(ctx context.Context, event *Event) error {
					data, err := MarshalEvent(event)
					require.NoError(err)

					ch <- data
					return nil
				},
			}
			ts := httptest.NewServer(wh)
			defer ts.Close()

			fuzz.Fuzz(&tc.payload)
			tc.header[signatureHeader] = "sha1=" +
				hex.EncodeToString(genMAC(tc.payload, wh.SecretKey, sha1.New))

			req, err := http.NewRequest(tc.method, ts.URL, bytes.NewBuffer(tc.payload))
			require.NoError(err)

			for k, v := range tc.header {
				req.Header.Set(k, v)
			}

			resp, err := ts.Client().Do(req)
			require.NoError(err)

			respBody, err := ioutil.ReadAll(resp.Body)
			require.NoError(err)
			resp.Body.Close()

			require.Truef(tc.reponseCode == resp.StatusCode,
				"status(%d): %s, expected(%d): %s, error: %s\n",
				resp.StatusCode, http.StatusText(resp.StatusCode),
				tc.reponseCode, http.StatusText(tc.reponseCode),
				respBody,
			)

			if len(tc.payload) > 0 && tc.reponseCode == http.StatusOK {
				data := <-ch
				event, err := UnmarshalEvent(data)
				require.NoError(err)

				require.Truef(event.Type == tc.header[eventTypeHeader],
					"event type: %s, expected: %s\n",
					event.Type, tc.header[eventTypeHeader],
				)

				require.Truef(bytes.Equal(event.Payload, tc.payload),
					"event payload: %s, expected: %s\n",
					event.Payload, tc.payload,
				)
			}
		})
	}
}

// genMAC generates the HMAC signature for a message provided the secret key
// and hashFunc.
func genMAC(message, key []byte, hashFunc func() hash.Hash) []byte {
	mac := hmac.New(hashFunc, key)
	mac.Write(message)
	return mac.Sum(nil)
}

// checkMAC reports whether messageMAC is a valid HMAC tag for message.
func checkMAC(message, messageMAC, key []byte, hashFunc func() hash.Hash) bool {
	expectedMAC := genMAC(message, key, hashFunc)
	return hmac.Equal(messageMAC, expectedMAC)
}

// messageMAC returns the hex-decoded HMAC tag from the signature and its
// corresponding hash function.
func messageMAC(signature string) ([]byte, func() hash.Hash, error) {
	if signature == "" {
		return nil, nil, errors.New("missing signature")
	}
	sigParts := strings.SplitN(signature, "=", 2)
	if len(sigParts) != 2 {
		return nil, nil, fmt.Errorf("error parsing signature %q", signature)
	}

	var hashFunc func() hash.Hash
	switch sigParts[0] {
	case sha1Prefix:
		hashFunc = sha1.New
	case sha256Prefix:
		hashFunc = sha256.New
	case sha512Prefix:
		hashFunc = sha512.New
	default:
		return nil, nil, fmt.Errorf("unknown hash type prefix: %q", sigParts[0])
	}

	buf, err := hex.DecodeString(sigParts[1])
	if err != nil {
		return nil, nil, fmt.Errorf("error decoding signature %q: %v", signature, err)
	}
	return buf, hashFunc, nil
}
