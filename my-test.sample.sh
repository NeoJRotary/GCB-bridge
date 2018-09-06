#!/bin/bash
# run test inside docker to simulate container runtime enviroment

# prepare ENV
GITHUB_APP_PRIVATE_KEY=$(cat private.pem)
GITHUB_APP_ID="app id"
GITHUB_APP_WEBHOOK_SECRET="webhook secret"
GCLOUD_SERVICE_ACCOUNT=$(cat service-account.json)
GCLOUD_PROJECT_ID="project id"
GCLOUD_STORAGE_BUCKET="storage bucket"

# switch test targets with 0 and 1
TEST_WEBHOOK=1
TEST_TRIGGER=1
TEST_GCLOUD=1
TEST_MAIN=1


if [ "$1" = "--rebuild-base" ]
then
  docker build -t neojrotary/gcb-bridge/test-base -f test-base.Dockerfile .
fi

docker build -t neojrotary/gcb-bridge/test -f test.Dockerfile .

docker run --rm \
  -e GITHUB_APP_PRIVATE_KEY="$GITHUB_APP_PRIVATE_KEY" \
  -e GITHUB_APP_ID="$GITHUB_APP_ID" \
  -e GITHUB_APP_WEBHOOK_SECRET="$GITHUB_APP_WEBHOOK_SECRET" \
  -e GCLOUD_SERVICE_ACCOUNT="$GCLOUD_SERVICE_ACCOUNT" \
  -e GCLOUD_PROJECT_ID="$GCLOUD_PROJECT_ID" \
  -e GCLOUD_STORAGE_BUCKET="$GCLOUD_STORAGE_BUCKET" \
  -e TEST_WEBHOOK="$TEST_WEBHOOK" \
  -e TEST_TRIGGER="$TEST_TRIGGER" \
  -e TEST_GCLOUD="$TEST_GCLOUD" \
  -e TEST_MAIN="$TEST_MAIN" \
  neojrotary/gcb-bridge/test
