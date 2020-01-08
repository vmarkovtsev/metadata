package metadata

import (
	"context"
	"os"
	"strconv"

	"github.com/athenianco/metadata/github"
	"github.com/athenianco/metadata/pubsub"
)

/*
 * Note: By default,
 * Cloud Functions does not support connecting to the Cloud SQL instance using TCP.
 * Your code should not try to access the instance using an IP address (such as 127.0.0.1 or 172.17.0.1)
 * unless you have configured Serverless VPC Access.
 *
 * The PostgreSQL standard requires the Unix socket to have a .s.PGSQL.5432 suffix in the socket path.
 * Some libraries apply this suffix automatically,
 * but others require you to specify the socket path as follows:
 * /cloudsql/INSTANCE_CONNECTION_NAME/.s.PGSQL.5432.
 */
var (
	githubProcessor pubsub.Subscriber
)

func init() {

	dbURI := os.Getenv("GITHUB_DATABASE_URI")
	if dbURI == "" {
		panic("GITHUB_DATABASE_URI is not set")
	}

	// When using a connection pool, it is important to set the maximum connections to 1.
	// This may seem counter-intuitive, however, Cloud Functions limits concurrent executions to 1 per instance.
	// This means you never have a situation where two requests are being processed by a single function instance at the same time.
	// This means in most situations only a single database connection is needed.
	maxOpenConns, err := strconv.Atoi(os.Getenv("GITHUB_DATABASE_MAX_OPEN_CONNS"))
	if err != nil {
		maxOpenConns = 1
	}

	maxIdleConns, err := strconv.Atoi(os.Getenv("GITHUB_DATABASE_MAX_IDLE_CONNS"))
	if err != nil {
		maxIdleConns = 1
	}

	db, err := github.OpenDatabase(dbURI, maxOpenConns, maxIdleConns)
	if err != nil {
		panic(err)
	}

	githubProcessor = func(ctx context.Context, msg pubsub.Message) error {
		event, err := github.UnmarshalEvent(msg.Data)
		if err != nil {
			return err
		}
		return event.Process(ctx, db)
	}
}

// GithubProcess is triggered by Pub/Sub.
func GithubProcess(ctx context.Context, msg pubsub.Message) error {
	return githubProcessor(ctx, msg)
}
