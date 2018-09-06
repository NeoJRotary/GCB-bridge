package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
	"testing"
	"time"

	"github.com/NeoJRotary/GCB-bridge/gcloud"
	"github.com/NeoJRotary/GCB-bridge/trigger"
	D "github.com/NeoJRotary/describe-go"
	"github.com/NeoJRotary/describe-go/dhttp"
	"github.com/NeoJRotary/exec-go"
)

// TestMain integration test of [ webhook > trigger > gcloud > github ]
// prepare "github-event-payload.json" and "github-event-headers.json" at project root

func TestMain(t *testing.T) {
	exec.DefaultEventHandler = &exec.EventHandler{
		CmdStarted: func(cmd *exec.Cmd) {
			fmt.Println("[ CMD ]", cmd.GetCmd())
		},
	}
	trigger.DEBUG = true

	go main()
	time.Sleep(time.Second * 2)

	// load sample json
	payloadB, err := ioutil.ReadFile("./github-event-payload.json")
	if D.IsErr(err) {
		t.Fatal(err)
	}

	buf := new(bytes.Buffer)
	err = json.Compact(buf, payloadB)
	if D.IsErr(err) {
		t.Fatal(err)
	}
	payloadB = buf.Bytes()

	headersB, err := ioutil.ReadFile("./github-event-headers.json")
	if D.IsErr(err) {
		t.Fatal(err)
	}

	var headers map[string]string
	err = json.Unmarshal(headersB, &headers)
	if D.IsErr(err) {
		t.Fatal(err)
	}

	var payload map[string]interface{}
	err = json.Unmarshal(payloadB, &payload)
	if D.IsErr(err) {
		t.Fatal(err)
	}

	// enable msg listener
	gcloud.InitMessageListener(strconv.FormatInt(time.Now().Unix(), 10))

	// send webhook to localhost server
	res, err := dhttp.Client(dhttp.TypeClient{
		Method: "POST",
		URL:    "http://localhost:8080/webhook",
		Header: headers,
		Body:   payloadB,
	}).Do()
	if D.IsErr(err) {
		t.Fatal(err)
	}
	if res.StatusCode != 200 {
		t.Fatal("should receive 200, get " + res.Status)
	}

	// waiting for gcloud event
	status := gcloud.MsgListener.ListenStatus(time.Second * 60)
	if status != "QUEUED" {
		t.Fatal("should receive QUEUED")
	}
	status = gcloud.MsgListener.ListenStatus(time.Second * 100)
	if status != "WORKING" {
		t.Fatal("should receive WORKING")
	}
	status = gcloud.MsgListener.ListenStatus(time.Second * 100)
	if status != "SUCCESS" {
		t.Fatal("should receive SUCCESS")
	}
}
