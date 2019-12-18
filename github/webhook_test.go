package github

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	fuzz "github.com/google/gofuzz"
)

func TestWebhook(t *testing.T) {
	tests := []struct {
		method      string
		header      map[string]string
		payload     []byte
		reponseCode int
	}{
		{
			method: "POST",
			header: map[string]string{
				SignatureHeader: "SECRET_TOKEN",
				EventTypeHeader: "TEST",
			},
			reponseCode: http.StatusOK,
		},
		{
			method:      "GET",
			header:      map[string]string{},
			reponseCode: http.StatusBadRequest,
		},
		{
			method:      "GET",
			header:      map[string]string{SignatureHeader: "SECRET_TOKEN"},
			reponseCode: http.StatusMethodNotAllowed,
		},
		{
			method:      "POST",
			header:      map[string]string{EventTypeHeader: "TEST"},
			reponseCode: http.StatusBadRequest,
		},
		{
			method:      "POST",
			header:      map[string]string{SignatureHeader: "SECRET_TOKEN"},
			reponseCode: http.StatusBadRequest,
		},
	}

	ch := make(chan []byte, 1)
	wh := NewWebhook(
		func(ctx context.Context, event *Event) error {
			data, err := MarshalEvent(event)
			if err != nil {
				return err
			}

			ch <- data
			return nil
		},
	)

	ts := httptest.NewServer(wh)
	defer ts.Close()

	fuzz := fuzz.New()
	for _, tc := range tests {
		fuzz.Fuzz(&tc.payload)

		req, err := http.NewRequest(tc.method, ts.URL, bytes.NewBuffer(tc.payload))
		if err != nil {
			t.Fatal(err)
		}
		for k, v := range tc.header {
			req.Header.Set(k, v)
		}

		resp, err := ts.Client().Do(req)
		if err != nil {
			t.Fatal(err)
		}

		respBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Fatal(err)
		}
		resp.Body.Close()

		if tc.reponseCode != resp.StatusCode {
			t.Fatalf("status(%d): %s, expected(%d): %s, error: %s\n",
				resp.StatusCode, http.StatusText(resp.StatusCode),
				tc.reponseCode, http.StatusText(tc.reponseCode),
				respBody,
			)
		}

		if len(tc.payload) > 0 && tc.reponseCode == http.StatusOK {
			data := <-ch
			event, err := UnmarshalEvent(data)
			if err != nil {
				t.Fatal(err)
			}

			if event.Type != tc.header[EventTypeHeader] {
				t.Errorf("event type: %s, expected: %s\n", event.Type, tc.header[EventTypeHeader])
			}

			if !bytes.Equal(event.Payload, tc.payload) {
				t.Errorf("event payload: %s, expected: %s\n", event.Payload, tc.payload)
			}
		}
	}
}
