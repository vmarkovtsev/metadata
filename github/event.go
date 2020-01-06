package github

import (
	"context"
	"encoding/json"

	gh "github.com/google/go-github/v28/github"
)

// Event is a Github's event received by webhooks.
type Event struct {
	// Event types: https://developer.github.com/v3/activity/events/types/
	Type string `json:"type"`

	// Payload should be json.RawMessage (for optimization and what's expected),
	// but for safety and fuzzy testing we use []byte,
	// because based on doc. "[]byte encodes as a base64-encoded string".
	//
	// delay parsing until we know the Type
	Payload []byte `json:"payload,omitempty"`
}

// UpsertPayload parses the event payload and upserts it to the given database.
func (e *Event) UpsertPayload(ctx context.Context, db *Database) error {
	event, err := gh.ParseWebHook(e.Type, e.Payload)
	if err != nil {
		return err
	}

	switch event := event.(type) {
	case *gh.RepositoryEvent:
		// Triggered when a repository is created, archived, unarchived, renamed, edited, transferred, enabled
		// for anonymous Git access, disabled for anonymous Git access, made public, or made private.
		return db.UpsertRepository(ctx, event)

	case *gh.OrganizationEvent:
		// Triggered when an organization is created, renamed, and deleted,
		// and when a user is added, removed, or invited to an organization.
		// Global webhooks will only receive notifications when an organization is created and deleted.
		// Organization webhooks will receive notifications for deleted, added, removed, renamed, and invited events.
		return db.UpsertOrganization(ctx, event)

	case *gh.IssueCommentEvent:
		// IssueCommentEvent is triggered when an issue comment is created on an issue
		// or pull request.
		return db.UpsertIssueComment(ctx, event)

	case *gh.IssuesEvent:
		// IssuesEvent is triggered when an issue is opened, edited, deleted, transferred,
		// pinned, unpinned, closed, reopened, assigned, unassigned, labeled, unlabeled,
		// locked, unlocked, milestoned, or demilestoned.
		return db.UpsertIssues(ctx, event)

	case *gh.PullRequestEvent:
		// Triggered when a pull request is assigned, unassigned, labeled, unlabeled,
		// opened, edited, closed, reopened, synchronize, ready_for_review,
		// locked, unlocked or when a pull request review is requested or removed.
		return db.UpsertPullRequest(ctx, event)

	case *gh.PullRequestReviewEvent:
		// PullRequestReviewEvent is triggered when a review is submitted on a pull
		// request.
		return db.UpsertPullRequestReview(ctx, event)

	case *gh.PullRequestReviewCommentEvent:
		// Triggered when a comment on a pull request's unified diff is created,
		// edited, or deleted (in the Files Changed tab).
		return db.UpsertPullRequestReviewComment(ctx, event)
	}

	return nil
}

// UnmarshalEvent parses data and save result as an Event.
func UnmarshalEvent(data []byte) (*Event, error) {
	var event Event
	if err := json.Unmarshal(data, &event); err != nil {
		return nil, err
	}
	return &event, nil
}

// MarshalEvent returns encoded Event.
func MarshalEvent(event *Event) ([]byte, error) {
	return json.Marshal(event)
}
