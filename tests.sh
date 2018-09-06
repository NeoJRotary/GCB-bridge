#!/bin/bash

# test sets be run inside test docker

if [ "$TEST_WEBHOOK" = "1" ]
then
  echo "TEST_WEBHOOK"
  go test -v ./webhook
fi

if [ "$TEST_TRIGGER" = "1" ]
then
  echo "TEST_TRIGGER"
  go test -v ./trigger
fi

if [ "$TEST_GCLOUD" = "1" ]
then
  echo "TEST_GCLOUD"
  go test -v ./gcloud
fi

if [ "$TEST_MAIN" = "1" ]
then
  echo "TEST_MAIN"
  go test -v -run Main
fi
