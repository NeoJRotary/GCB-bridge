#!/bin/bash
REPO_NAME="neojrotary/gcb-bridge"

if [ "$1" == "" ]
then
  docker build --rm -t $REPO_NAME:latest .
  docker push $REPO_NAME:latest
else
  docker build --rm -t $REPO_NAME:latest -t $REPO_NAME:$1 .
  docker push $REPO_NAME:latest
  docker push $REPO_NAME:$1
fi
