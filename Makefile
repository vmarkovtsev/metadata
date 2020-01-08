RUNTIME = go111
REGION ?= "us-central1"

# Note: By default,
# Cloud Functions does not support connecting to the Cloud SQL instance using TCP.
# Your code should not try to access the instance using an IP address (such as 127.0.0.1 or 172.17.0.1)
# unless you have configured Serverless VPC Access.
#
# The PostgreSQL standard requires the Unix socket to have a .s.PGSQL.5432 suffix in the socket path.
# Some libraries apply this suffix automatically,
# but others require you to specify the socket path as follows:
# postgres://<db_user>:<db_pass>@/<db_name>?host=/cloudsql/<cloud_sql_instance_name>/.s.PGSQL.5432
GITHUB_DATABASE_URI ?= "postgres://user:password@127.0.0.1:5432/metadata"

# When using a connection pool, it is important to set the maximum connections to 1.
# This may seem counter-intuitive, however, Cloud Functions limits concurrent executions to 1 per instance.
# This means you never have a situation where two requests are being processed by a single function instance at the same time.
# This means in most situations only a single database connection is needed.
GITHUB_DATABASE_MAX_OPEN_CONNS ?= 1
GITHUB_DATABASE_MAX_IDLE_CONNS ?= 1

GITHUB_WEBHOOK_TOPIC ?= "github-hook-owl"
GITHUB_WEBHOOK_SECRET_KEY ?= "secret-token"

GITHUB_WEBHOOK_NAME = "github-hook-owl"
GITHUB_WEBHOOK_ENTRY_POINT = "GithubWebhook"

GITHUB_PROCESSOR_NAME = "github_processor"
GITHUB_PROCESSOR_ENTRY_POINT = "GithubProcess"

echo-vars:
	@echo \
	"RUNTIME=${RUNTIME}\n"\
	"REGION=${REGION}\n"\
	"GITHUB_DATABASE_URI=${GITHUB_DATABASE_URI}\n"\
	"GITHUB_WEBHOOK_TOPIC=${GITHUB_WEBHOOK_TOPIC}\n"\
	"GITHUB_WEBHOOK_NAME=${GITHUB_WEBHOOK_NAME}\n"\
	"GITHUB_WEBHOOK_ENTRY_POINT=${GITHUB_WEBHOOK_ENTRY_POINT}\n"\
	"GITHUB_WEBHOOK_SECRET_KEY=${GITHUB_WEBHOOK_SECRET_KEY}\n"\
	"GITHUB_PROCESSOR_NAME=${GITHUB_PROCESSOR_NAME}\n"\
	"GITHUB_PROCESSOR_ENTRY_POINT=${GITHUB_PROCESSOR_ENTRY_POINT}\n"\
	"GITHUB_DATABASE_MAX_OPEN_CONNS=${GITHUB_DATABASE_MAX_OPEN_CONNS}\n"\
	"GITHUB_DATABASE_MAX_IDLE_CONNS=${GITHUB_DATABASE_MAX_IDLE_CONNS}\n"


test-all:
	# cloud_sql_proxy -instances=<cloud_sql_instance_name>=tcp:5432
	go test -v ./... --count 1

create-github-webhook-topic:
	gcloud pubsub topics create $(GITHUB_WEBHOOK_TOPIC) --message-storage-policy-allowed-regions $(REGION)

deploy-github-processor:
	gcloud functions deploy $(GITHUB_PROCESSOR_NAME) --entry-point $(GITHUB_PROCESSOR_ENTRY_POINT) \
	--trigger-topic $(GITHUB_WEBHOOK_TOPIC) \
	--runtime $(RUNTIME) \
	--region $(REGION) \
	--set-env-vars GITHUB_DATABASE_URI=$(GITHUB_DATABASE_URI) \
	--ignore-file ".gcloudignore"

deploy-github-webhook:
	gcloud functions deploy $(GITHUB_WEBHOOK_NAME) --entry-point $(GITHUB_WEBHOOK_ENTRY_POINT) \
	--trigger-http \
	--runtime $(RUNTIME) \
	--region $(REGION) \
	--set-env-vars GITHUB_WEBHOOK_TOPIC=$(GITHUB_WEBHOOK_TOPIC) --set-env-vars GITHUB_WEBHOOK_SECRET_KEY=$(GITHUB_WEBHOOK_SECRET_KEY) \
	--ignore-file ".gcloudignore"

deploy-all:	create-github-webhook-topic	deploy-github-processor	deploy-github-webhook
