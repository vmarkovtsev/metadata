package github

import (
	"context"
	"net/http"

	gh "github.com/google/go-github/v28/github"
)

// Webhook is the implementation of http.Handler with OnEvent callback.
type Webhook struct {
	SecretKey []byte
	OnEvent   func(ctx context.Context, event *Event) error
}

func (h *Webhook) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	payload, err := gh.ValidatePayload(r, h.SecretKey)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	typ := gh.WebHookType(r)
	if typ == "" {
		http.Error(w, "invalid webhook event type", http.StatusBadRequest)
		return
	}

	err = h.OnEvent(r.Context(),
		&Event{
			Type:    typ,
			Payload: payload,
		},
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
