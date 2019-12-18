package github

import (
	"context"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

const (
	// SignatureHeader is a hash signature passed along with each request.
	SignatureHeader = "X-Hub-Signature"

	// EventTypeHeader is a request header used to know which event has been received.
	EventTypeHeader = "X-GitHub-Event"

	maxPayloadSize = 1024 * 1024
)

// Webhook is the implementation of http.Handler with OnEvent callback.
type Webhook struct {
	Handler http.Handler
	OnEvent func(ctx context.Context, event *Event) error
}

// NewWebhook creates a new instance with default validators
func NewWebhook(onevent func(ctx context.Context, event *Event) error) *Webhook {
	wh := Webhook{OnEvent: onevent}

	// validators are nested, so will be called from the last one
	// up to the top (the most nested one).
	var h http.Handler = http.HandlerFunc(wh.serveHTTP)
	h = eventTypeValidator{h}
	h = methodValidator{h}
	h = signatureValidator{h}

	wh.Handler = h
	return &wh
}

func (wh *Webhook) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	wh.Handler.ServeHTTP(w, r)
}

func (wh *Webhook) serveHTTP(w http.ResponseWriter, r *http.Request) {
	if wh.OnEvent == nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	payload, err := ioutil.ReadAll(io.LimitReader(r.Body, maxPayloadSize))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = wh.OnEvent(r.Context(),
		&Event{
			Type:    r.Header.Get(EventTypeHeader),
			Payload: payload,
		},
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// signatureValidator is http.Handler which validates signature header.
type signatureValidator struct {
	h http.Handler
}

func (h signatureValidator) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	signature := r.Header.Get(SignatureHeader)
	if signature == "" {
		http.Error(w, "invalid signature: "+signature, http.StatusBadRequest)
		return
	}
	if h.h != nil {
		h.h.ServeHTTP(w, r)
	}
}

// methodValidator is http.Handler which validates request method.
type methodValidator struct {
	h http.Handler
}

func (h methodValidator) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch strings.ToUpper(r.Method) {
	case "POST":
		if h.h != nil {
			h.h.ServeHTTP(w, r)
		}

	default:
		http.Error(w, "invalid request method: "+r.Method, http.StatusMethodNotAllowed)
	}
}

// eventTypeValidator is http.Handler which validates event type header.
type eventTypeValidator struct {
	h http.Handler
}

func (h eventTypeValidator) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	eventType := r.Header.Get(EventTypeHeader)
	if eventType == "" {
		http.Error(w, "invalid event type: "+eventType, http.StatusBadRequest)
		return
	}
	if h.h != nil {
		h.h.ServeHTTP(w, r)
	}
}
