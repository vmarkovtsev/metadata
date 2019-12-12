package pubsub

import (
	"context"

	gcloud "cloud.google.com/go/pubsub"
)

const (
	projectID = "athenian-1"
	topicID   = "metadata"
)

func Publish(ctx context.Context, payload []byte) error {
	cli, err := gcloud.NewClient(ctx, projectID)
	if err != nil {
		return err
	}

	topic := cli.Topic(topicID)
	msg := &gcloud.Message{
		Data: payload,
	}

	_, err = topic.Publish(ctx, msg).Get(ctx)
	if err != nil {
		return err
	}

	return nil
}
