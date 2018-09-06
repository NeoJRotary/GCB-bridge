package webhook

import (
	"encoding/json"
	"log"
	"strconv"

	"github.com/NeoJRotary/GCB-bridge/app"
	"github.com/NeoJRotary/GCB-bridge/github"
	"github.com/NeoJRotary/GCB-bridge/trigger"
	D "github.com/NeoJRotary/describe-go"
)

type pullRequest struct {
	Head pullRequestRef `json:"head"`
	Base pullRequestRef `json:"base"`
}

type pullRequestEvent struct {
	Action       string       `json:"action"`
	Number       int          `json:"number"`
	PullRequest  pullRequest  `json:"pull_request"`
	Repository   repository   `json:"repository"`
	Installation installation `json:"installation"`
}

type pushEvent struct {
	Ref          string       `json:"ref"`
	Before       string       `json:"before"`
	After        string       `json:"after"`
	Created      bool         `json:"created"`
	Deleted      bool         `json:"deleted"`
	Repository   repository   `json:"repository"`
	Installation installation `json:"installation"`
}

type repository struct {
	ID       int    `json:"id"`
	NodeID   string `json:"node_id"`
	FullName string `json:"full_name"`
}

type installation struct {
	ID int `json:"id"`
}

type pullRequestRef struct {
	Ref string `json:"ref"`
	SHA string `json:"sha"`
}

func eventHandler(evtName string, body []byte) {
	switch evtName {
	case "push":
		pushEventHandler(body)
	case "pull_request":
		pullRequestEventHandler(body)
	default:
	}
}

func pushEventHandler(body []byte) {
	var payload pushEvent
	err := json.Unmarshal(body, &payload)
	if D.IsErr(err) {
		log.Println("PushEvent Unmarshal Error", err)
		return
	}

	// dont handle delete
	if payload.Deleted {
		return
	}

	refD := D.String(payload.Ref)
	switch true {
	// push to branch
	case refD.HasPrefix("refs/heads/"):
		repo := app.Repo{
			Event:            "Branch",
			InstallationID:   strconv.Itoa(payload.Installation.ID),
			RepositoryNodeID: payload.Repository.NodeID,
			FullName:         payload.Repository.FullName,
			Branch:           refD.TrimPrefix("refs/heads/").Get(),
			BeforeSHA:        payload.Before,
			AfterSHA:         payload.After,
		}

		PRList, err := github.GetAssociatedPullRequests(repo.InstallationID, repo.FullName, repo.Branch)
		if !D.IsErr(err) {
			ss := D.StringSlice()
			for _, pr := range PRList {
				ss.Push(pr.BaseRefName)
			}
			repo.AssociatedBases = ss.Get()
		}

		trigger.EventHandler(&repo)
	case refD.HasPrefix("refs/tags/"):
		repo := app.Repo{
			Event:            "Tag",
			InstallationID:   strconv.Itoa(payload.Installation.ID),
			RepositoryNodeID: payload.Repository.NodeID,
			FullName:         payload.Repository.FullName,
			Tag:              refD.TrimPrefix("refs/tags/").Get(),
			BeforeSHA:        payload.Before,
			AfterSHA:         payload.After,
		}
		trigger.EventHandler(&repo)
	default:
		log.Println("Invalid Ref", payload.Ref)
	}
}

func pullRequestEventHandler(body []byte) {
	var payload pullRequestEvent
	err := json.Unmarshal(body, &payload)
	if D.IsErr(err) {
		log.Println("PullRequestEven Unmarshal Error", err)
		return
	}

	// only handle PR opened
	if payload.Action != "opened" {
		return
	}

	repo := app.Repo{
		Event:            "PullRequest",
		InstallationID:   strconv.Itoa(payload.Installation.ID),
		RepositoryNodeID: payload.Repository.NodeID,
		FullName:         payload.Repository.FullName,
		Branch:           payload.PullRequest.Head.Ref,
		BaseBranch:       payload.PullRequest.Base.Ref,
		BeforeSHA:        payload.PullRequest.Base.SHA,
		AfterSHA:         payload.PullRequest.Head.SHA,
	}
	trigger.EventHandler(&repo)
}
