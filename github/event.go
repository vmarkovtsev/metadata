package github

import (
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

// ParsePayload parses the event payload. For recognized event types.
func (e *Event) ParsePayload() error {
	switch e.Type {
	// Triggered when a commit comment is created.
	case "commit_comment":
		var payload gh.CommitCommentEvent
		return json.Unmarshal(e.Payload, &payload)

	// Represents a created branch or tag.
	case "create":
		var payload gh.CreateEvent
		return json.Unmarshal(e.Payload, &payload)

	// Represents a deleted branch or tag.
	case "delete":
		var payload gh.DeleteEvent
		return json.Unmarshal(e.Payload, &payload)

	// Triggered when an issue comment is created, edited, or deleted.
	case "issue_comment":
		var payload gh.IssueCommentEvent
		return json.Unmarshal(e.Payload, &payload)

	// Triggered when an issue is opened, edited, deleted, pinned, unpinned,
	// closed, reopened, assigned, unassigned, labeled, unlabeled,
	// locked, unlocked, transferred, milestoned, or demilestoned.
	case "issues":
		var payload gh.IssuesEvent
		return json.Unmarshal(e.Payload, &payload)

	// Triggered when the webhook that this event is configured on is deleted.
	// This event will only listen for changes to the particular hook the event is installed on.
	// Therefore, it must be selected for each hook that you'd like to recieve meta events for.
	case "meta":
		var payload gh.MetaEvent
		return json.Unmarshal(e.Payload, &payload)

	// Triggered when an organization is created, renamed, and deleted,
	// and when a user is added, removed, or invited to an organization.
	// Global webhooks will only receive notifications when an organization is created and deleted.
	// Organization webhooks will receive notifications for deleted, added, removed, renamed, and invited events.
	case "organization":
		var payload gh.OrganizationEvent
		return json.Unmarshal(e.Payload, &payload)

	// Triggered when a project is created, updated, closed, reopened, or deleted.
	case "project":
		var payload gh.ProjectEvent
		return json.Unmarshal(e.Payload, &payload)

	// Triggered when a private repository is made public.
	// Without a doubt: the best GitHub Enterprise Server event.
	case "public":
		var payload gh.PublicEvent
		return json.Unmarshal(e.Payload, &payload)

	// Triggered when a pull request is assigned, unassigned, labeled, unlabeled,
	// opened, edited, closed, reopened, synchronize, ready_for_review,
	// locked, unlocked or when a pull request review is requested or removed.
	case "pull_request":
		var payload gh.PullRequestEvent
		return json.Unmarshal(e.Payload, &payload)

	// Triggered when a pull request review is submitted into a non-pending state,
	// the body is edited, or the review is dismissed.
	case "pull_request_review":
		var payload gh.PullRequestReviewEvent
		return json.Unmarshal(e.Payload, &payload)

	// Triggered when a comment on a pull request's unified diff is created, edited, or deleted (in the Files Changed tab).
	case "pull_request_review_comment":
		var payload gh.PullRequestReviewCommentEvent
		return json.Unmarshal(e.Payload, &payload)

	// Triggered on a push to a repository branch.
	// Branch pushes and repository tag pushes also trigger webhook push events.
	case "push":
		var payload gh.PushEvent
		return json.Unmarshal(e.Payload, &payload)

	// Triggered when a release is published, unpublished, created, edited, deleted, or prereleased.
	case "release":
		var payload gh.ReleaseEvent
		return json.Unmarshal(e.Payload, &payload)

	// Triggered when a repository is created, archived, unarchived, renamed, edited, transferred, enabled
	// for anonymous Git access, disabled for anonymous Git access, made public, or made private.
	// Organization hooks are also triggered when a repository is deleted.
	case "repository":
		var payload gh.RepositoryEvent
		return json.Unmarshal(e.Payload, &payload)

	// Triggered when the status of a Git commit changes.
	case "status":
		var payload gh.StatusEvent
		return json.Unmarshal(e.Payload, &payload)
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
