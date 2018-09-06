# GCB-Bridge v0.1.0
CI integration for Github App and Google Cloud Build

[Docker Hub](https://hub.docker.com/r/neojrotary/gcb-bridge/)   
`docker pull neojrotary/gcb-bridge`   
   
## Introduction
### Concept
We want to use amazing GCB with Github CheckRun but...
- GCB Github App cant manage triggers
- GCB doesn't support triggers inside configuration yaml   

before Google support them, we make a small app to handle it ( at 2018/09/06). **This Project will be deprecated when Google support them.**

### Features 
- Integration with Github CheckRun on commits.
- Check Build Log at CheckRun detail without access GCB.
- Multi-Builds in one configuration yaml.
- Define triggers inside configuration yaml.
- Triggers can be set at both `build` and `step`.
- Trigger by `branch`, `tag`, `pull request` and filtering by file changes.

### How To Use
Prepare `cloudbuild.bridge.yaml` at root of repository. You can
- Define `name` for each build, it will be showed at Github CheckRun.
- Separate builds by `---`.
- Put `triggers` at build and step.
- Each build follow "Cloud Build build configuration". Check [Build Configuration Overview](https://cloud.google.com/cloud-build/docs/build-config)

For example: 
```
name: 'test branch'
triggers:
- branches: []
steps:
- name: ''
  triggers:
  - includedFiles: []
- name: ''
---
name: 'new release'
triggers:
- tags: []
steps:
- name: ''
- name: ''
images:
- [...]
artifacts:
```
   
`triggers` is an array of `trigger Object`, a `triggers` is passed if any of `trigger Object` is passed. Each `trigger Object` takes same concept of "GCB Build Triggers".
```
triggers:
- branches: []
  tags: []
  pullRequestBases: []
  includedFiles: []
  ignoredFiles: []
- ... trigger Object...
- ... trigger Object...
```
- `branches` : array of regex   
  Match branch by regex.
- `tags` : array of regex   
  Match tag by regex.
- `pullRequestBases` : array of regex   
  Match Pull-Request's base branch by regex.
- `includedFiles` : array of glob   
  Match file changes by glob. It works by itself even there is no any `branches`, `tag`, `pullRequestBases`.
- `ignoredFiles`: array of glob   
  Ignore file changes by glob.

## Installation
### How To
- Prepare your private Github App. Check [Creating a GitHub App](https://developer.github.com/apps/building-github-apps/creating-a-github-app/). Since it is a private App, the only URL you need to setup correctly is `Webhook URL`. Server will listen on `/webook` so it would be like `https://sub.domain.com/webhook`.
- Prepare your Google Cloud Platform Project and create a Service Account with `Cloud Build Service Account`, `Pub/Sub Subscriber` and `Storage Admin` roles. Get credential file in json format.
- Deploy container with ENVs to any Engine you like in GCP. Then setup network for receiving Github Webhook.

### ENV
Enviroment variables for runtime
- `LISTEN_ON`   
  Address where http server listen on. Default is `0.0.0.0:8080`.
- `GITHUB_APP_PRIVATE_KEY`   
  Generated key file from App's settings. Put the content into ENV.
- `GITHUB_APP_ID`   
  Get it from App's settings "About" section.
- `GITHUB_APP_WEBHOOK_SECRET`   
  Your Github app webhook secret
- `GCLOUD_SERVICE_ACCOUNT`   
  GCP service account json file content
- `GCLOUD_PROJECT_ID`   
  GCP Project ID
- `GCLOUD_STORAGE_BUCKET`   
  GCS bucket name, your log and source file will be uploaded to here
- `DEBUG`   
  Enable DEBUG mode to print more information to console

## Develop
If you want to know deep into project. Check [DEVELOP.md](https://github.com/NeoJRotary/GCB-bridge/blob/master/DEVELOP.md).

## What's Next
- improve logging
- optional integration with Slack
- support `waitFor` between different builds
- support auto cancellation if new build be triggered but previous same build still not finish
- prepare REST API for manully triggering
