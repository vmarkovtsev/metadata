RUNTIME = go111
REGION ?= "us-central1"

GITHUB_WEBHOOK_TOPIC ?= "github-hook-owl"

GITHUB_WEBHOOK_NAME = "github_webhook"
GITHUB_WEBHOOK_ENTRY_POINT = "GithubWebhook"

GITHUB_UPDATE_NAME = "github_update"
GITHUB_UPDATE_ENTRY_POINT = "GithubUpdate"

echo-vars:
	@echo "RUNTIME=${RUNTIME}\nREGION=${REGION}\nGITHUB_WEBHOOK_TOPIC=${GITHUB_WEBHOOK_TOPIC}\nGITHUB_WEBHOOK_NAME=${GITHUB_WEBHOOK_NAME}\nGITHUB_WEBHOOK_ENTRY_POINT=${GITHUB_WEBHOOK_ENTRY_POINT}\nGITHUB_UPDATE_NAME=${GITHUB_UPDATE_NAME}\nGITHUB_UPDATE_ENTRY_POINT=${GITHUB_UPDATE_ENTRY_POINT}\n"

test-all:
	go test -v ./... --count 1

deploy-github-update:
	gcloud functions deploy $(GITHUB_UPDATE_NAME) --entry-point $(GITHUB_UPDATE_ENTRY_POINT) --trigger-topic $(GITHUB_WEBHOOK_TOPIC) --runtime $(RUNTIME) --region $(REGION) --ignore-file "${GITHUB_WEBHOOK_NAME}.gcloudignore"

deploy-github-webhook:
	gcloud functions deploy $(GITHUB_WEBHOOK_NAME) --entry-point $(GITHUB_WEBHOOK_ENTRY_POINT) --trigger-http --runtime $(RUNTIME) --region $(REGION) --set-env-vars GITHUB_WEBHOOK_TOPIC=$(GITHUB_WEBHOOK_TOPIC) --ignore-file "${GITHUB_WEBHOOK_NAME}.gcloudignore"

deploy-all: deploy-github-update	deploy-github-webhook
