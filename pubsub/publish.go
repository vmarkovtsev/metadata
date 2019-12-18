package pubsub

import (
	"context"
	"log"
	"os"

	gcp "cloud.google.com/go/pubsub"
)

// projectID is set from the GCP_PROJECT environment variable, which is
// automatically set by the Cloud Functions runtime.
var projectID = os.Getenv("GCP_PROJECT")

// Publisher is Google Pub/Sub publisher.
type Publisher struct {
	topic *gcp.Topic
}

// NewPublisher creates a new instance of Pub/Sub publisher.
// It also creates the Pub/Sub topic if it does not exist.
func NewPublisher(topicID string) (*Publisher, error) {
	ctx := context.Background()

	client, err := gcp.NewClient(ctx, projectID)
	if err != nil {
		return nil, err
	}

	// Create the topic if it doesn't exist.
	topic := client.Topic(topicID)
	exists, err := topic.Exists(ctx)
	if err != nil {
		return nil, err
	}
	if !exists {
		if _, err = client.CreateTopic(ctx, topicID); err != nil {
			return nil, err
		}
	}

	return &Publisher{topic: topic}, nil
}

// Publish data to the Pub/Sub topic synchronously.
func (p *Publisher) Publish(ctx context.Context, data []byte) error {
	id, err := p.topic.Publish(ctx, &gcp.Message{Data: data}).Get(ctx)
	if err != nil {
		log.Printf("publish id: %s, data: %s, error: %v\n", id, string(data), err)
	}
	return err
}
