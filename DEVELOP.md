# Develop
Happy to see you are interesting in developing of this project :)

## Concect
Reform `cloudbuild.bridge.yaml` to `cloudbuild.yaml` for each valid build and parallelly submit them by `gcloud builds submit`. 
   
Webhook Handling Flow :
- Check event type
- Clone repo and check `cloudbuild.bridge.yaml` exists or not
- Read file then Split Builds
- Validate builds
- Collect valid builds
- Reform valid builds and save to temp file
- Submit builds parallelly
- Listen on PubSub
- Update CheckRun to `QUEUED` when build is `QUEUED`
- Update CheckRun to `IN_PROGESS` when build is `WORKING`
- Update CheckRun to `COMPLETED` when build is finished.

## Testing at local
You will need `private.key` and `service-account.json` at project root. Prepare your own `my-test.sh` and execute it. Adjust $ENV at `my-test.sh` to control which part need to test. Check [my-test.sample.sh](https://github.com/NeoJRotary/GCB-bridge/blob/master/my-test.sample.sh) and [tests.sh](https://github.com/NeoJRotary/GCB-bridge/blob/master/tests.sh) for detail.   
Some tests may need `github-event-headers.json` and `github-event-payload.json` to simulate contents of Github Webhook.