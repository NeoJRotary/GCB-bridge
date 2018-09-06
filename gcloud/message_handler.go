package gcloud

import (
	"encoding/json"
	"log"
	"strconv"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/NeoJRotary/GCB-bridge/github"
	D "github.com/NeoJRotary/describe-go"
	exec "github.com/NeoJRotary/exec-go"
)

// PubSubPrintMessageData pubsub print message data
var PubSubPrintMessageData = false

type buildMsg struct {
	ID            string `json:"id"`
	Status        string `json:"status"`
	LogURL        string `json:"logUrl"`
	CreateTime    string `json:"createTime"`
	StartTime     string `json:"startTime"`
	FinishTime    string `json:"finishTime"`
	Substitutions struct {
		Branch           string `json:"BRANCH_NAME"`
		CommitSHA        string `json:"COMMIT_SHA"`
		RepoName         string `json:"REPO_NAME"`
		Tag              string `json:"TAG_NAME"`
		InstallationID   string `json:"_GITHUB_INSTALLATION_ID"`
		RepositoryNodeID string `json:"_GITHUB_REPOSITORY_NODE_ID"`
		CheckRunID       string `json:"_GITHUB_CHECKRUN_ID"`
		BridgeUID        string `json:"_BRIDGE_UID"`
	} `json:"substitutions"`
	Timing struct {
		BUILD       timing `json:"BUILD"`
		FETCHSOURCE timing `json:"FETCHSOURCE"`
	} `json:"timing"`
}

type timing struct {
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
}

func messageHandler(m *pubsub.Message) {
	if PubSubPrintMessageData {
		log.Printf("----- PubSub Message -----\n" + string(m.Data) + "\n----------\n")
	}

	var msg buildMsg
	err := json.Unmarshal(m.Data, &msg)
	if D.IsErr(err) {
		log.Println("PubSub messageHandler Unmarshal error", err)
	}

	if EnableMsgListener {
		MsgListener.receive(&msg)
	}

	// stop calling github API
	if SkipGithubAPI {
		return
	}

	// send to github
	switch msg.Status {
	case "QUEUED":
		err = github.QueuedCheckRun(
			msg.Substitutions.InstallationID,
			msg.Substitutions.RepositoryNodeID,
			msg.Substitutions.CheckRunID,
			msg.ID,
			msg.LogURL,
		)
	case "WORKING":
		err = github.InProgressCheckRun(
			msg.Substitutions.InstallationID,
			msg.Substitutions.RepositoryNodeID,
			msg.Substitutions.CheckRunID,
		)
	case "CANCELLED":
		out := getBuildLog(&msg)
		err = github.CompletedCheckRun(
			msg.Substitutions.InstallationID,
			msg.Substitutions.RepositoryNodeID,
			msg.Substitutions.CheckRunID,
			github.CheckCancelled,
			msg.FinishTime,
			getSummary(&msg),
			out,
		)
	case "TIMEOUT":
		out := getBuildLog(&msg)
		err = github.CompletedCheckRun(
			msg.Substitutions.InstallationID,
			msg.Substitutions.RepositoryNodeID,
			msg.Substitutions.CheckRunID,
			github.CheckTimeout,
			msg.FinishTime,
			getSummary(&msg),
			out,
		)
	case "SUCCESS":
		out := getBuildLog(&msg)
		err = github.CompletedCheckRun(
			msg.Substitutions.InstallationID,
			msg.Substitutions.RepositoryNodeID,
			msg.Substitutions.CheckRunID,
			github.CheckSuccess,
			msg.FinishTime,
			getSummary(&msg),
			out,
		)
	case "FAILURE":
		out := getBuildLog(&msg)
		err = github.CompletedCheckRun(
			msg.Substitutions.InstallationID,
			msg.Substitutions.RepositoryNodeID,
			msg.Substitutions.CheckRunID,
			github.CheckFailure,
			msg.FinishTime,
			getSummary(&msg),
			out,
		)
	}

	if D.IsErr(err) {
		log.Println("PubSub messageHandler error", msg.Status, err)
	}
}

func getBuildDuration(msg *buildMsg) string {
	start, err := time.Parse(time.RFC3339, msg.StartTime)
	if D.IsErr(err) {
		log.Println("PubSub messageHandler getBuildDuration error", err)
	}
	finish, err := time.Parse(time.RFC3339, msg.FinishTime)
	if D.IsErr(err) {
		log.Println("PubSub messageHandler getBuildDuration error", err)
	}
	dur := finish.Sub(start)
	return strconv.FormatFloat(dur.Seconds(), 'f', 5, 64)
}

func getSummary(msg *buildMsg) string {
	return `Status : ` + msg.Status + `
	Started : ` + msg.StartTime + `
	Build time : ` + getBuildDuration(msg) + ` sec`
}

func getBuildLog(msg *buildMsg) string {
	out, err := exec.RunCmd("", "gcloud", "builds", "log", msg.ID)
	if D.IsErr(err) {
		log.Println("PubSub messageHandler gcloud log error", err)
		return ""
	}

	return "```\n" + out + "```"
}
