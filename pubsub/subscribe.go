package pubsub

import (
	"context"
)

// Message is the payload of a Pub/Sub event.
type Message struct {
	Data []byte `json:"data"`
}

// Subscriber is Pub/Sub push subscriber.
type Subscriber func(ctx context.Context, msg Message) error
