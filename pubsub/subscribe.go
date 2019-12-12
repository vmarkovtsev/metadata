package pubsub

import (
	"context"

	gcloud "cloud.google.com/go/pubsub"
)

func Subscribe(ctx context.Context, msg gcloud.Message) error {
	return nil
}
