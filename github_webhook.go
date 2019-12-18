package metadata

import (
	"context"
	"net/http"
	"os"

	"github.com/athenianco/metadata/github"
	"github.com/athenianco/metadata/pubsub"
)

var githubWebhook *github.Webhook

func init() {
	topicID := os.Getenv("GITHUB_WEBHOOK_TOPIC")
	if topicID == "" {
		panic("GITHUB_WEBHOOK_TOPIC is not set")
	}

	publisher, err := pubsub.NewPublisher(topicID)
	if err != nil {
		panic(err)
	}

	githubWebhook = github.NewWebhook(
		func(ctx context.Context, event *github.Event) error {
			data, err := github.MarshalEvent(event)
			if err != nil {
				return err
			}
			return publisher.Publish(ctx, data)
		},
	)
}

// GithubWebhook is http.Handler triggered by github on metadata events.
func GithubWebhook(w http.ResponseWriter, r *http.Request) {
	githubWebhook.ServeHTTP(w, r)
}
