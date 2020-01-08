package github

import (
	"context"
	"encoding/json"
	"fmt"

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

// Process parses the event payload and upserts/alters it to the given database.
func (e *Event) Process(ctx context.Context, db *Database) error {
	event, err := gh.ParseWebHook(e.Type, e.Payload)
	if err != nil {
		return err
	}

	switch event := event.(type) {
	case *gh.InstallationEvent:
		// Triggered when someone installs (created) , uninstalls (deleted),
		// or accepts new permissions (new_permissions_accepted) for a GitHub App.
		return processInstallationEvent(ctx, db, event)

	case *gh.InstallationRepositoriesEvent:
		// Triggered when a repository is added or removed from an installation.
		return processInstallationRepositoriesEvent(ctx, db, event)

	case *gh.RepositoryEvent:
		// Triggered when a repository is created, archived, unarchived, renamed, edited, transferred, enabled
		// for anonymous Git access, disabled for anonymous Git access, made public, or made private.
		return processRepositoryEvent(ctx, db, event)

	case *gh.OrganizationEvent:
		// Triggered when an organization is created, renamed, and deleted,
		// and when a user is added, removed, or invited to an organization.
		// Global webhooks will only receive notifications when an organization is created and deleted.
		// Organization webhooks will receive notifications for deleted, added, removed, renamed, and invited events.
		return processOrganizationEvent(ctx, db, event)

	case *gh.IssueCommentEvent:
		// IssueCommentEvent is triggered when an issue comment is created on an issue
		// or pull request.
		return processIssueCommentEvent(ctx, db, event)

	case *gh.IssuesEvent:
		// IssuesEvent is triggered when an issue is opened, edited, deleted, transferred,
		// pinned, unpinned, closed, reopened, assigned, unassigned, labeled, unlabeled,
		// locked, unlocked, milestoned, or demilestoned.
		return processIssuesEvent(ctx, db, event)

	case *gh.PullRequestEvent:
		// Triggered when a pull request is assigned, unassigned, labeled, unlabeled,
		// opened, edited, closed, reopened, synchronize, ready_for_review,
		// locked, unlocked or when a pull request review is requested or removed.
		return processPullRequestEvent(ctx, db, event)

	case *gh.PullRequestReviewEvent:
		// PullRequestReviewEvent is triggered when a review is submitted on a pull
		// request.
		return processPullRequestReviewEvent(ctx, db, event)

	case *gh.PullRequestReviewCommentEvent:
		// Triggered when a comment on a pull request's unified diff is created,
		// edited, or deleted (in the Files Changed tab).
		return processPullRequestReviewCommentEvent(ctx, db, event)
	}

	return nil
}

func processInstallationEvent(ctx context.Context, db *Database, event *gh.InstallationEvent) (err error) {
	defer errRecover(event, &err)

	switch event.GetAction() {
	case "deleted":
		break

	case "created":
		// TODO(@kuba--): Probably, somehow we should check
		// if event.GetInstallation().GetAppID() is our App!
		// And ignore other installations.
		for _, repo := range event.Repositories {
			err = db.UpsertRepository(ctx, repo)
			if err != nil {
				break
			}
		}
	}

	return err
}

func processInstallationRepositoriesEvent(ctx context.Context, db *Database, event *gh.InstallationRepositoriesEvent) (err error) {
	defer errRecover(event, &err)

	switch event.GetAction() {
	case "removed":
		break

	case "added":
		for _, repo := range event.RepositoriesAdded {
			err = db.UpsertRepository(ctx, repo)
			if err != nil {
				break
			}
		}
	}

	return err
}

func processRepositoryEvent(ctx context.Context, db *Database, event *gh.RepositoryEvent) (err error) {
	defer errRecover(event, &err)

	switch event.GetAction() {
	case "deleted":
		break

	case "anonymous_access_enabled", "anonymous_access_disabled":
		break

	case "created", "edited", "renamed", "transferred", "archived", "unarchived", "publicized", "privatized":
		return db.UpsertRepository(ctx, event.GetRepo())
	}

	return err
}

func processOrganizationEvent(ctx context.Context, db *Database, event *gh.OrganizationEvent) (err error) {
	defer errRecover(event, &err)

	switch event.GetAction() {
	case "deleted", "member_removed":
		break

	case "created", "renamed", "member_added", "member_invited":
		return db.UpsertOrganization(ctx, event.GetOrganization())
	}

	return err
}

func processIssueCommentEvent(ctx context.Context, db *Database, event *gh.IssueCommentEvent) (err error) {
	defer errRecover(event, &err)

	switch event.GetAction() {
	case "deleted":
		break

	case "created", "edited":
		if event.GetIssue().IsPullRequest() {
			err = db.UpsertIssueCommentAsPullRequest(ctx, event.GetRepo(), event.GetIssue(), event.GetComment())
		} else {
			err = db.UpsertIssueComment(ctx, event.GetRepo(), event.GetIssue(), event.GetComment())
		}
	}

	return err
}

func processIssuesEvent(ctx context.Context, db *Database, event *gh.IssuesEvent) (err error) {
	defer errRecover(event, &err)

	switch event.GetAction() {
	case "deleted", "unpinned", "unassigned", "unlocked", "demilestoned":
		break

	case "opened",
		"edited",
		"reopened",
		"transferred",
		"labeled",
		"assigned",
		"closed",
		"locked",
		"pinned",
		"milestoned":
		return db.UpsertIssues(ctx, event.GetRepo(), event.GetIssue())
	}

	return err
}

func processPullRequestEvent(ctx context.Context, db *Database, event *gh.PullRequestEvent) (err error) {
	defer errRecover(event, &err)

	switch event.GetAction() {
	case "closed", "synchronize", "locked", "unlocked":
		break

	case "assigned",
		"unassigned",
		"labeled",
		"unlabeled",
		"opened",
		"edited",
		"reopened",
		"ready_for_review":
		return db.UpsertPullRequest(ctx, event.GetRepo(), event.GetPullRequest())
	}

	return err
}

func processPullRequestReviewEvent(ctx context.Context, db *Database, event *gh.PullRequestReviewEvent) (err error) {
	defer errRecover(event, &err)

	switch event.GetAction() {
	case "dismissed":
		break

	case "submitted", "edited":
		return db.UpsertPullRequestReview(ctx, event.GetRepo(), event.GetPullRequest(), event.GetReview())
	}

	return err
}

func processPullRequestReviewCommentEvent(ctx context.Context, db *Database, event *gh.PullRequestReviewCommentEvent) (err error) {
	defer errRecover(event, &err)

	switch event.GetAction() {
	case "deleted":
		break

	case "created", "edited":
		return db.UpsertPullRequestReviewComment(ctx, event.GetRepo(), event.GetPullRequest(), event.GetComment())
	}

	return err
}

func errRecover(event interface{}, err *error) {
	if r := recover(); r != nil {
		*err = fmt.Errorf("Event(%s) recovered from: %v", gh.Stringify(event), r)
	}
}
