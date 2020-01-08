package metadata

import (
	"context"
	"net/http"
	"os"
	"sync"

	"github.com/athenianco/metadata/github"
	"github.com/athenianco/metadata/pubsub"
)

var ghWebhook struct {
	once sync.Once
	*github.Webhook
}

func initGHWebhook() {
	topicID := os.Getenv("GITHUB_WEBHOOK_TOPIC")
	if topicID == "" {
		panic("GITHUB_WEBHOOK_TOPIC is not set")
	}

	secretKey := os.Getenv("GITHUB_WEBHOOK_SECRET_KEY")
	if secretKey == "" {
		panic("GITHUB_WEBHOOK_SECRET_KEY is not set")
	}

	publisher, err := pubsub.NewPublisher(topicID)
	if err != nil {
		panic(err)
	}

	ghWebhook.Webhook = &github.Webhook{
		SecretKey: []byte(secretKey),
		OnEvent: func(ctx context.Context, event *github.Event) error {
			data, err := github.MarshalEvent(event)
			if err != nil {
				return err
			}
			return publisher.Publish(ctx, data)
		},
	}
}

// GithubWebhook is http.Handler triggered by github on metadata events.
func GithubWebhook(w http.ResponseWriter, r *http.Request) {
	ghWebhook.once.Do(initGHWebhook)
	ghWebhook.ServeHTTP(w, r)
}
