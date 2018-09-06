package github

import (
	"encoding/json"
	"errors"

	"github.com/NeoJRotary/GCB-bridge/app"
	D "github.com/NeoJRotary/describe-go"
	"github.com/NeoJRotary/describe-go/dhttp"
)

// GraphqlRes graphql response
type GraphqlRes struct {
	Data   GraphqlData    `json:"data"`
	Errors []GraphqlError `json:"errors"`
}

// GraphqlError graphql error
type GraphqlError struct {
	Message string `json:"message"`
}

// GraphqlData graphql data
type GraphqlData struct {
	Repository struct {
		Ref struct {
			AssociatedPullRequests struct {
				Nodes []PullRequest `json:"nodes"`
			} `json:"associatedPullRequests"`
		} `json:"ref"`

		PullRequests struct {
			Nodes []PullRequest `json:"nodes"`
		} `json:"pullRequests"`
	} `json:"repository"`

	CreateCheckRun struct {
		ClientMutationID string   `json:"clientMutationId"`
		CheckRun         CheckRun `json:"checkRun"`
	} `json:"createCheckRun"`

	UpdateCheckRun struct {
		ClientMutationID string   `json:"clientMutationId"`
		CheckRun         CheckRun `json:"checkRun"`
	} `json:"updateCheckRun"`
}

func apiPOST(installationID, query string, variables map[string]interface{}) (*GraphqlData, error) {
	b, _ := json.Marshal(map[string]interface{}{
		"query":     query,
		"variables": variables,
	})

	res, err := dhttp.Client(dhttp.TypeClient{
		Method: "POST",
		URL:    "https://api.github.com/graphql",
		Header: map[string]string{
			"Accept":        "application/vnd.github.antiope-preview+json",
			"Authorization": "Bearer " + app.GetAccessToken(installationID),
		},
		Body: b,
	}).Do()

	if D.IsErr(err) {
		return nil, err
	}

	// not 2xx
	if res.StatusCode >= 300 {
		b := res.ReadAllBody()
		return nil, D.NewErr("Github apiPOST", res.Status, string(b))
	}

	// unmarshal body
	var gres GraphqlRes
	err = json.Unmarshal(res.ReadAllBody(), &gres)
	if D.IsErr(err) {
		return nil, err
	}

	if len(gres.Errors) != 0 {
		msgs := D.StringSlice()
		for _, e := range gres.Errors {
			msgs.Push(e.Message)
		}
		return nil, errors.New(msgs.Join(" ; ").Get())
	}

	return &gres.Data, nil
}
