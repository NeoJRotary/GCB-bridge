package webhook

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"testing"

	D "github.com/NeoJRotary/describe-go"
)

// prepare "github-event-payload.json" and "github-event-headers.json" at project root

func TestWebhook(t *testing.T) {
	payloadB, err := ioutil.ReadFile("../github-event-payload.json")
	if D.IsErr(err) {
		t.Fatal(err)
	}

	buf := new(bytes.Buffer)
	err = json.Compact(buf, payloadB)
	if D.IsErr(err) {
		t.Fatal(err)
	}
	payloadB = buf.Bytes()

	headersB, err := ioutil.ReadFile("../github-event-headers.json")
	if D.IsErr(err) {
		t.Fatal(err)
	}

	var headers map[string]string
	err = json.Unmarshal(headersB, &headers)
	if D.IsErr(err) {
		t.Fatal(err)
	}

	if !verifyWebhook(payloadB, headers["X-Hub-Signature"]) {
		t.Fatal("should get true")
	}
}
