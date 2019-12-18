package metadata

import (
	"context"

	"github.com/athenianco/metadata/github"
	"github.com/athenianco/metadata/pubsub"
)

var githubDbUpdater pubsub.Subscriber

func init() {
	githubDbUpdater = func(ctx context.Context, msg pubsub.Message) error {
		_, err := github.UnmarshalEvent(msg.Data)
		if err != nil {
			return err
		}

		// not implemented yet.
		return nil
	}
}

// GithubUpdate is triggered by Pub/Sub.
func GithubUpdate(ctx context.Context, msg pubsub.Message) error {
	return githubDbUpdater(ctx, msg)
}
