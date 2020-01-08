package github

import (
	"bytes"
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"log"
	"time"

	gh "github.com/google/go-github/v28/github"
	"github.com/lib/pq"
)

var tables = map[string]string{
	"github_organizations_versioned":         "avatar_url, collaborators, created_at, description, email, htmlurl, id, login, name, node_id, owned_private_repos, public_repos, total_private_repos, updated_at",
	"github_users_versioned":                 "avatar_url, bio, company, created_at, email, followers, following, hireable, htmlurl, id, location, login, name, node_id, organization_id, organization_login, owned_private_repos, private_gists, public_gists, public_repos, total_private_repos, updated_at",
	"github_repositories_versioned":          "allow_merge_commit, allow_rebase_merge, allow_squash_merge, archived, created_at, default_branch, description, disabled, fork, forks_count, fullname, has_issues, has_wiki, homepage, htmlurl, id, language, name, node_id, open_issues_count, owner_id, owner_login, owner_type, private, pushed_at, sshurl, stargazers_count, topics, updated_at, watchers_count",
	"github_issues_versioned":                "assignees, body, closed_at, closed_by_id, closed_by_login, comments, created_at, htmlurl, id, labels, locked, milestone_id, milestone_title, node_id, number, repository_name, repository_owner, repository_fullname, state, title, updated_at, user_id, user_login",
	"github_issue_comments_versioned":        "author_association, body, created_at, htmlurl, id, issue_number, node_id, repository_name, repository_owner, repository_fullname, updated_at, user_id, user_login",
	"github_pull_requests_versioned":         "additions, assignees, author_association, base_ref, base_repository_name, base_repository_owner, base_repository_fullname, base_sha, base_user, body, changed_files, closed_at, comments, commits, created_at, deletions, head_ref, head_repository_name, head_repository_owner, head_repository_fullname, head_sha, head_user, htmlurl, id, labels, maintainer_can_modify, merge_commit_sha, mergeable, merged, merged_at, merged_by_id, merged_by_login, milestone_id, milestone_title, node_id, number, repository_name, repository_owner, repository_fullname, review_comments, state, title, updated_at, user_id, user_login",
	"github_pull_request_reviews_versioned":  "body, commit_id, htmlurl, id, node_id, pull_request_number, repository_name, repository_owner, repository_fullname, state, submitted_at, user_id, user_login",
	"github_pull_request_comments_versioned": "author_association, body, commit_id, created_at, diff_hunk, htmlurl, id, in_reply_to, node_id, original_commit_id, original_position, path, position, pull_request_number, pull_request_review_id, repository_name, repository_owner, repository_fullname, updated_at, user_id, user_login",
}

// Database is a postgres database where github metadata are stored.
type Database struct {
	*sql.DB
}

// OpenDatabase opens postgres connection
func OpenDatabase(dbURI string, maxOpenConns, maxIdleConns int) (*Database, error) {
	db, err := sql.Open("postgres", dbURI)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(maxOpenConns)
	db.SetMaxIdleConns(maxIdleConns)

	if err = db.Ping(); err != nil {
		db.Close()
		return nil, err
	}
	return &Database{db}, nil
}

func (db *Database) txExecContext(ctx context.Context, query string, args ...interface{}) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	if res, err := tx.ExecContext(ctx, query, args...); err != nil {
		tx.Rollback()
		log.Printf("query: %s, args: %v, result: %s, error: %v\n", query, args, gh.Stringify(res), err)
		return err
	}
	return tx.Commit()
}

// UpsertRepository (github_repositories_versioned)
func (db *Database) UpsertRepository(ctx context.Context, repo *gh.Repository) error {
	const tab = "github_repositories_versioned"
	cols := tables[tab]
	ver := version()

	topics := make([]string, len(repo.Topics))
	for i, t := range repo.Topics {
		topics[i] = t
	}
	query := fmt.Sprintf(`
	INSERT INTO %s
	(sum256, versions, %s)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14,
		$15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29,
		$30, $31, $32)
	ON CONFLICT (sum256)
	DO UPDATE
	SET versions = array_append(%s.versions, $33)`, tab, cols, tab)
	return db.txExecContext(ctx, query,
		sum256(repo.GetID()),       // sum256,
		pq.Array([]int64{ver}),     // versions,
		repo.GetAllowMergeCommit(), // allow_merge_commit boolean
		repo.GetAllowRebaseMerge(), // allow_rebase_merge boolean
		repo.GetAllowSquashMerge(), // allow_squash_merge boolean
		repo.GetArchived(),         // archived boolean
		repo.GetCreatedAt().UTC(),  // created_at timestamptz
		repo.GetDefaultBranch(),    // default_branch text
		repo.GetDescription(),      // description text
		repo.GetDisabled(),         // disabled boolean
		repo.GetFork(),             // fork boolean
		repo.GetForksCount(),       // forks_count bigint
		repo.GetFullName(),         // fullname text
		repo.GetHasIssues(),        // has_issues boolean
		repo.GetHasWiki(),          // has_wiki boolean
		repo.GetHomepage(),         // homepage text
		repo.GetHTMLURL(),          // htmlurl text
		repo.GetID(),               // id bigint,
		repo.GetLanguage(),         // language text
		repo.GetName(),             // name text
		repo.GetNodeID(),           // node_id text
		repo.GetOpenIssuesCount(),  // open_issues_count bigint
		repo.GetOwner().GetID(),    // owner_id bigint NOT NULL,
		repo.GetOwner().GetLogin(), // owner_login text NOT NULL,
		repo.GetOwner().GetType(),  // owner_type text NOT NULL
		repo.GetPrivate(),          // private boolean
		repo.GetPushedAt().UTC(),   // pushed_at timestamptz
		repo.GetSSHURL(),           // sshurl text
		repo.GetStargazersCount(),  // stargazers_count bigint
		pq.Array(topics),           // topics text[] NOT NULL
		repo.GetUpdatedAt().UTC(),  // updated_at timestamptz
		repo.GetWatchersCount(),    // watchers_count bigint
		ver,
	)
}

// UpsertOrganization (github_organizations_versioned)
func (db *Database) UpsertOrganization(ctx context.Context, org *gh.Organization) error {
	const tab = "github_organizations_versioned"
	cols := tables[tab]
	ver := version()

	query := fmt.Sprintf(`
		INSERT INTO %s
		(sum256, versions, %s)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14,
			$15, $16)
		ON CONFLICT (sum256)
		DO UPDATE
		SET versions = array_append(%s.versions, $17)`, tab, cols, tab)
	return db.txExecContext(ctx, query,
		sum256(org.GetID()),        // sum256,
		pq.Array([]int64{ver}),     // versions,
		org.GetAvatarURL(),         // avatar_url text,
		org.GetCollaborators(),     // collaborators bigint,
		org.GetCreatedAt(),         // created_at timestamptz,
		org.GetDescription(),       // description text,
		org.GetEmail(),             // email text,
		org.GetHTMLURL(),           // htmlurl text,
		org.GetID(),                // id bigint,
		org.GetLogin(),             // login text,
		org.GetName(),              // name text,
		org.GetNodeID(),            // node_id text,
		org.GetOwnedPrivateRepos(), // owned_private_repos bigint,
		org.GetPublicRepos(),       // public_repos bigint,
		org.GetTotalPrivateRepos(), // total_private_repos bigint,
		org.GetUpdatedAt(),         // updated_at timestamptz,
		ver,
	)
}

// UpsertPullRequest (github_pull_requests_versioned)
func (db *Database) UpsertPullRequest(ctx context.Context, repo *gh.Repository, pr *gh.PullRequest) error {
	const tab = "github_pull_requests_versioned"
	cols := tables[tab]
	ver := version()

	var assignees []string = make([]string, len(pr.Assignees))
	for i, a := range pr.Assignees {
		assignees[i] = a.GetLogin()
	}

	var labels []string = make([]string, len(pr.Labels))
	for i, l := range pr.Labels {
		labels[i] = l.GetName()
	}

	query := fmt.Sprintf(
		`INSERT INTO %s
		(sum256, versions, %s)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14,
			$15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29,
			$30, $31, $32, $33, $34, $35, $36, $37, $38, $39, $40, $41, $42, $43, $44,
			$45, $46, $47)
		ON CONFLICT (sum256)
		DO UPDATE
		SET versions = array_append(%s.versions, $48)`, tab, cols, tab)
	return db.txExecContext(ctx, query,
		sum256(
			repo.GetID(),
			pr.GetID(),
		), // sum256,
		pq.Array([]int64{ver}),                       // versions,
		pr.GetAdditions(),                            // additions bigint,
		pq.Array(assignees),                          // assignees text[] NOT NULL,
		pr.GetAuthorAssociation(),                    // author_association text,
		pr.GetBase().GetRef(),                        // base_ref text NOT NULL,
		pr.Base.GetRepo().GetName(),                  // base_repository_name text NOT NULL,
		pr.GetBase().GetRepo().GetOwner().GetLogin(), // base_repository_owner text NOT NULL,
		pr.GetBase().GetRepo().GetFullName(),         // base_repository_fullname
		pr.GetBase().GetSHA(),                        // base_sha text NOT NULL,
		pr.GetBase().GetUser().GetLogin(),            // base_user text NOT NULL,
		pr.GetBody(),                                 // body text,
		pr.GetChangedFiles(),                         // changed_files bigint,
		pr.GetClosedAt(),                             // closed_at timestamptz,
		pr.GetComments(),                             // comments bigint,
		pr.GetCommits(),                              // commits bigint,
		pr.GetCreatedAt(),                            // created_at timestamptz,
		pr.GetDeletions(),                            // deletions bigint,
		pr.GetHead().GetRef(),                        // head_ref text NOT NULL,
		pr.GetHead().GetRepo().GetName(),             // head_repository_name text NOT NULL,
		pr.GetHead().GetRepo().GetOwner().GetLogin(), // head_repository_owner text NOT NULL,
		pr.GetHead().GetRepo().GetFullName(),         // head_repository_fullname
		pr.GetHead().GetSHA(),                        // head_sha text NOT NULL,
		pr.GetHead().GetUser().GetLogin(),            // head_user text NOT NULL,
		pr.GetHTMLURL(),                              // htmlurl text,
		pr.GetID(),                                   // id bigint,
		pq.Array(labels),                             // labels text[] NOT NULL,
		pr.GetMaintainerCanModify(),                  // maintainer_can_modify boolean,
		pr.GetMergeCommitSHA(),                       // merge_commit_sha text,
		pr.GetMergeable(),                            // mergeable boolean,
		pr.GetMerged(),                               // merged boolean,
		pr.GetMergedAt(),                             // merged_at timestamptz,
		pr.GetMergedBy().GetID(),                     // merged_by_id bigint NOT NULL,
		pr.GetMergedBy().GetLogin(),                  // merged_by_login text NOT NULL,
		pr.GetMilestone().GetID(),                    // milestone_id text NOT NULL,
		pr.GetMilestone().GetTitle(),                 // milestone_title text NOT NULL,
		pr.GetID(),                                   // node_id text,
		pr.GetNumber(),                               // number bigint,
		repo.GetName(),                               // repository_name text NOT NULL,
		repo.GetOwner().GetLogin(),                   // repository_owner text NOT NULL,
		repo.GetFullName(),                           // repository_fullname
		pr.GetReviewComments(),                       // review_comments bigint,
		pr.GetState(),                                // state text,
		pr.GetTitle(),                                // title text,
		pr.GetUpdatedAt(),                            // updated_at timestamptz,
		pr.GetUser().GetID(),                         // user_id bigint NOT NULL,
		pr.GetUser().GetLogin(),                      // user_login text NOT NULL,
		ver,
	)
}

// UpsertPullRequestReviewComment (github_pull_request_comments_versioned)
func (db *Database) UpsertPullRequestReviewComment(ctx context.Context, repo *gh.Repository, pr *gh.PullRequest, comment *gh.PullRequestComment) error {
	const tab = "github_pull_request_comments_versioned"
	cols := tables[tab]
	ver := version()

	query := fmt.Sprintf(`
	INSERT INTO %s
	(sum256, versions, %s)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23)
	ON CONFLICT (sum256)
	DO UPDATE
	SET versions = array_append(%s.versions, $24)`, tab, cols, tab)
	return db.txExecContext(ctx, query,
		sum256(
			repo.GetID(),
			pr.GetID(),
			comment.GetID(),
		), // sum256,
		pq.Array([]int64{ver}),         // versions,
		comment.GetAuthorAssociation(), // author_association text,
		comment.GetBody(),              // body text,
		comment.GetCommitID(),          // commit_id text,
		comment.GetCreatedAt(),         // created_at timestamptz,
		comment.GetDiffHunk(),          // diff_hunk text,
		comment.GetHTMLURL(),           // htmlurl text,
		comment.GetID(),                // id bigint,
		comment.GetInReplyTo(),         // in_reply_to bigint,
		comment.GetNodeID(),            // node_id text,
		comment.GetOriginalCommitID(),  // original_commit_id text,
		comment.GetOriginalPosition(),  // original_position bigint,
		comment.GetPath(),              // path text,
		comment.GetPosition(),          // position bigint,
		pr.GetNumber(),                 // pull_request_number bigint NOT NULL,
		pr.GetID(),                     // pull_request_review_id bigint,
		repo.GetName(),                 // repository_name text NOT NULL,
		repo.GetOwner().GetLogin(),     // repository_owner text NOT NULL,
		repo.GetFullName(),             // repository_fullname
		comment.GetUpdatedAt(),         // updated_at timestamptz,
		comment.GetUser().GetID(),      // user_id bigint NOT NULL,
		comment.GetUser().GetLogin(),   // user_login text NOT NULL,
		ver,
	)
}

// UpsertPullRequestReview (github_pull_request_reviews_versioned)
func (db *Database) UpsertPullRequestReview(ctx context.Context, repo *gh.Repository, pr *gh.PullRequest, review *gh.PullRequestReview) error {
	const tab = "github_pull_request_reviews_versioned"
	cols := tables[tab]
	ver := version()

	query := fmt.Sprintf(`
	INSERT INTO %s
	(sum256, versions, %s)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
	ON CONFLICT (sum256)
	DO UPDATE
	SET versions = array_append(%s.versions, $16)`, tab, cols, tab)
	return db.txExecContext(ctx, query,
		sum256(
			repo.GetID(),
			pr.GetID(),
			review.GetID(),
		), // sum256,
		pq.Array([]int64{ver}),      // versions,
		review.GetBody(),            // body text,
		review.GetCommitID(),        // commit_id text,
		review.GetHTMLURL(),         // htmlurl text,
		review.GetID(),              // id bigint,
		review.GetNodeID(),          // node_id text,
		pr.GetNumber(),              // pull_request_number bigint NOT NULL,
		repo.GetName(),              // repository_name text NOT NULL,
		repo.GetOwner().GetLogin(),  // repository_owner text NOT NULL,
		repo.GetFullName(),          // repository_fullname
		review.GetState(),           // state text,
		review.GetSubmittedAt(),     // submitted_at timestamptz,
		review.GetUser().GetID(),    // user_id bigint NOT NULL,
		review.GetUser().GetLogin(), // user_login text NOT NULL,
		ver,
	)
}

// UpsertIssues (github_issues_versioned)
func (db *Database) UpsertIssues(ctx context.Context, repo *gh.Repository, issue *gh.Issue) error {
	const tab = "github_issues_versioned"
	cols := tables[tab]
	ver := version()

	var assignees []string = make([]string, len(issue.Assignees))
	for i, a := range issue.Assignees {
		assignees[i] = a.GetLogin()
	}

	var labels []string = make([]string, len(issue.Labels))
	for i, l := range issue.Labels {
		labels[i] = l.GetName()
	}

	query := fmt.Sprintf(`
	INSERT INTO %s (sum256, versions, %s)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14,
		$15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25)
	ON CONFLICT (sum256)
	DO UPDATE
	SET versions = array_append(%s.versions, $26)`, tab, cols, tab)
	return db.txExecContext(ctx, query,
		sum256(
			repo.GetID(),
			issue.GetID(),
		), // sum256,
		pq.Array([]int64{ver}),          // versions,
		pq.Array(assignees),             // assignees text[] NOT NULL,
		issue.GetBody(),                 // body text,
		issue.GetClosedAt(),             // closed_at timestamptz,
		issue.GetClosedBy().GetID(),     // closed_by_id bigint NOT NULL
		issue.GetClosedBy().GetLogin(),  // closed_by_login text NOT NULL,
		issue.GetComments(),             // comments bigint,
		issue.GetCreatedAt(),            // created_at timestamptz,
		issue.GetURL(),                  // htmlurl text,
		issue.GetID(),                   // id bigint,
		pq.Array(labels),                // labels text[] NOT NULL,
		issue.GetLocked(),               // locked boolean,
		issue.GetMilestone().GetID(),    // milestone_id text NOT NULL,
		issue.GetMilestone().GetTitle(), // milestone_title text NOT NULL,
		issue.GetNodeID(),               // node_id text,
		issue.GetNumber(),               // number bigint,
		repo.GetName(),                  // repository_name text NOT NULL,
		repo.GetOwner().GetLogin(),      // repository_owner text NOT NULL,
		repo.GetFullName(),              // repository_fullname text NOT NULL,
		issue.GetState(),                // state text,
		issue.GetTitle(),                // title text,
		issue.GetUpdatedAt(),            // updated_at timestamptz,
		issue.GetUser().GetID(),         // user_id bigint NOT NULL,
		issue.GetUser().GetLogin(),      // user_login text NOT NULL,
		ver,
	)
}

// UpsertIssueComment (github_issue_comments_versioned)
func (db *Database) UpsertIssueComment(ctx context.Context, repo *gh.Repository, issue *gh.Issue, comment *gh.IssueComment) error {
	const tab = "github_issue_comments_versioned"
	cols := tables[tab]
	ver := version()

	query := fmt.Sprintf(`
	INSERT INTO %s (sum256, versions, %s)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
	ON CONFLICT (sum256)
	DO UPDATE
	SET versions = array_append(%s.versions, $16)`, tab, cols, tab)
	return db.txExecContext(ctx, query,
		sum256(
			repo.GetID(),
			issue.GetID(),
			comment.GetID(),
		), // sum256,
		pq.Array([]int64{ver}),         // versions,
		comment.GetAuthorAssociation(), // author_association,
		comment.GetBody(),              // body,
		comment.GetCreatedAt(),         // created_at,
		comment.GetHTMLURL(),           // htmlurl,
		comment.GetID(),                // id,
		issue.GetNumber(),              // issue_number,
		comment.GetNodeID(),            // node_id,
		repo.GetName(),                 // repository_name,
		repo.GetOwner().GetLogin(),     // repository_owner,
		repo.GetFullName(),             // repository_fullname text NOT NULL,
		comment.GetUpdatedAt(),         // updated_at,
		comment.GetUser().GetID(),      // user_id,
		comment.GetUser().GetLogin(),   // user_login,
		ver,
	)
}

// UpsertIssueCommentAsPullRequest (github_issue_comments_versioned)
func (db *Database) UpsertIssueCommentAsPullRequest(ctx context.Context, repo *gh.Repository, issue *gh.Issue, comment *gh.IssueComment) error {
	const tab = "github_pull_request_comments_versioned"
	cols := tables[tab]
	ver := version()

	query := fmt.Sprintf(`
		INSERT INTO %s
		(sum256, versions, %s)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23)
		ON CONFLICT (sum256)
		DO UPDATE
		SET versions = array_append(%s.versions, $24)`, tab, cols, tab)
	return db.txExecContext(ctx, query,
		sum256(
			repo.GetID(),
			issue.GetID(),
			comment.GetID(),
		), // sum256,
		pq.Array([]int64{ver}),       // versions,
		issue.IsPullRequest(),        // author_association text,
		comment.GetBody(),            // body text,
		"",                           // commit_id text,
		comment.GetCreatedAt(),       // created_at timestamptz,
		"",                           // diff_hunk text,
		comment.GetHTMLURL(),         // htmlurl text,
		comment.GetID(),              // id bigint,
		0,                            // in_reply_to bigint,
		comment.GetNodeID(),          // node_id text,
		"",                           // original_commit_id text,
		0,                            // original_position bigint,
		"",                           // path text,
		"",                           // position bigint,
		issue.GetNumber(),            // pull_request_number bigint NOT NULL,
		comment.GetID(),              // pull_request_review_id bigint,
		repo.GetName(),               // repository_name text NOT NULL,
		repo.GetOwner().GetLogin(),   // repository_owner text NOT NULL,
		repo.GetFullName(),           // repository_fullname text NOT NULL,
		comment.GetUpdatedAt(),       // updated_at timestamptz,
		comment.GetUser().GetID(),    // user_id bigint NOT NULL,
		comment.GetUser().GetLogin(), // user_login text NOT NULL,
		ver,
	)
}

func sum256(ids ...int64) string {
	buf := new(bytes.Buffer)
	hash := sha256.New()
	for _, id := range ids {
		binary.Write(buf, binary.LittleEndian, id)
		hash.Write(buf.Bytes())
		buf.Reset()
	}
	return hex.EncodeToString(hash.Sum(nil))
}

func version() int64 {
	return time.Now().UTC().Unix()
}
