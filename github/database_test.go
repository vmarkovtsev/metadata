package github

import (
	"context"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

const dbURI = "postgres://user:password@127.0.0.1:5432/test?sslmode=disable"

func TestProcess(t *testing.T) {
	require := require.New(t)

	tests := []struct {
		name     string
		fixture  string
		query    string
		expected interface{}
		err      bool
	}{
		{
			name:    "installation",
			fixture: "testdata/installation_event.json",
			query: `select id from github_repositories_versioned where (
				name='cuckoo' and
				fork=false
			)`,
			expected: []interface{}{int64(85718512)},
		},
		{
			name:    "repository",
			fixture: "testdata/repository_event.json",
			query: `select 1 from github_repositories_versioned where (
				id=118 and
				language='Ruby' and
				name='Hello-World' and
				fork=false and
				default_branch='master' and
				owner_login='Codertocat'
			)`,
			expected: []interface{}{int64(1)},
		},
		{
			name:    "organization",
			fixture: "testdata/organization_event.json",
			query: `select 1 from github_organizations_versioned where (
				id=6 and
				login='Octocoders' and
				node_id='MDEyOk9yZ2FuaXphdGlvbjY=' and
				collaborators=0
			)`,
			expected: []interface{}{int64(1)},
		},
		{
			name:    "issue_comment",
			fixture: "testdata/issue_comment_event.json",
			query: `select body from github_issue_comments_versioned where (
				id=2 and
				issue_number=1 and
				user_login='Codertocat'
			)`,
			expected: []interface{}{"You are totally right! I'll get this fixed right away."},
		},
		{
			name:    "issues",
			fixture: "testdata/issues_event.json",
			query: `select repository_name from github_issues_versioned where (
				id=10 and
				assignees='{Codertocat}' and
				number=1
			)`,
			expected: []interface{}{"Hello-World"},
		},
		{
			name:    "pull_request",
			fixture: "testdata/pull_request_event.json",
			query: `select title from github_pull_requests_versioned where (
				number=2 and
				base_sha='78a96099c3f442d7f6e8d1a7d07090091993e65a' and
				head_sha='14977a7b5485400124827221a04bfb474bcd72d1'
				)`,
			expected: []interface{}{"Update the README with new information."},
		},
		{
			name:    "pull_request_review",
			fixture: "testdata/pull_request_review_event.json",
			query: `select id from github_pull_request_reviews_versioned where (
				commit_id='14977a7b5485400124827221a04bfb474bcd72d1' and
				user_login='Codertocat'
				)`,
			expected: []interface{}{int64(2)},
		},
		{
			name:    "pull_request_review_comment",
			fixture: "testdata/pull_request_review_comment_event.json",
			query: `select diff_hunk from github_pull_request_comments_versioned where (
				commit_id='14977a7b5485400124827221a04bfb474bcd72d1' and
				path='README.md' and
				user_login='Codertocat'
			)`,
			expected: []interface{}{"@@ -1 +1 @@\n-# Hello-World"},
		},
		{
			name:    "repository",
			fixture: "testdata/empty_event.json",
			err:     true,
		},
		{
			name:    "ignore",
			fixture: "testdata/empty_event.json",
			err:     false,
		},
	}

	db, err := OpenDatabase(dbURI, 0, 0)
	require.NoError(err)
	defer db.Close()

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			f, err := os.Open(tc.fixture)
			require.NoError(err)
			defer f.Close()

			payload, err := ioutil.ReadAll(f)
			require.NoError(err)

			event := &Event{
				Type:    tc.name,
				Payload: payload,
			}

			err = event.Process(context.TODO(), db)
			if tc.err {
				require.Error(err)
			} else if tc.query != "" {
				require.NoError(err)

				rows, err := db.Query(tc.query)
				require.NoError(err)

				vals := make([]interface{}, 0)
				for rows.Next() {
					var v interface{}
					err = rows.Scan(&v)
					require.NoError(err)
					vals = append(vals, v)
				}
				require.NoError(rows.Close())
				require.Equal(tc.expected, vals)
			}
		})
	}
}
