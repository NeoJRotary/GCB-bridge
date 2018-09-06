package gcloud

import (
	"fmt"
	"path"
	"strconv"
	"testing"
	"time"

	"github.com/NeoJRotary/GCB-bridge/app"
	"github.com/NeoJRotary/exec-go"
)

func TestGCloud(t *testing.T) {
	exec.DefaultEventHandler = &exec.EventHandler{
		CmdStarted: func(cmd *exec.Cmd) {
			fmt.Println("[ CMD ]", cmd.GetCmd())
		},
	}

	Init()
	app.InitRepo()

	SkipGithubAPI = true
	// PubSubPrintMessageData = true

	repo := &app.Repo{
		FullName: "NeoJRotary/GCB-bridge-test",
		Branch:   "master",
	}
	if !repo.Init() {
		t.Fatal("repo init failed")
	}

	InitMessageListener(strconv.FormatInt(time.Now().Unix(), 10))

	go StartBuild(repo, "test-build", path.Join(repo.Dir, "test.success.yaml"))

	status := MsgListener.ListenStatus(time.Second * 20)
	if status != "QUEUED" {
		t.Fatal("should receive QUEUED")
	}
	status = MsgListener.ListenStatus(time.Second * 60)
	if status != "WORKING" {
		t.Fatal("should receive WORKING")
	}
	status = MsgListener.ListenStatus(time.Second * 60)
	if status != "SUCCESS" {
		t.Fatal("should receive SUCCESS")
	}

	go StartBuild(repo, "test-build", path.Join(repo.Dir, "test.fail.yaml"))

	status = MsgListener.ListenStatus(time.Second * 20)
	if status != "QUEUED" {
		t.Fatal("should receive QUEUED")
	}
	status = MsgListener.ListenStatus(time.Second * 60)
	if status != "WORKING" {
		t.Fatal("should receive WORKING")
	}
	status = MsgListener.ListenStatus(time.Second * 60)
	if status != "FAILURE" {
		t.Fatal("should receive FAILURE")
	}

}
