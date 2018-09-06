package github

import (
	"log"

	D "github.com/NeoJRotary/describe-go"
)

// CheckRun Github CheckRun
type CheckRun struct {
	ID string `json:"id"`
}

// CheckConclusionState CheckConclusionState
type CheckConclusionState string

const (
	// CheckCancelled CheckConclusionState CANCELLED
	CheckCancelled CheckConclusionState = "CANCELLED"
	// CheckFailure CheckConclusionState FAILURE
	CheckFailure CheckConclusionState = "FAILURE"
	// CheckSuccess CheckConclusionState SUCCESS
	CheckSuccess CheckConclusionState = "SUCCESS"
	// CheckTimeout CheckConclusionState TIMED_OUT
	CheckTimeout CheckConclusionState = "TIMED_OUT"
)

// ActionRequiredCheckRun create ACTION_REQUIRED CheckRun for invalid build
func ActionRequiredCheckRun(installationID, repositoryNodeID, headSHA, name, text string) {
	_, err := apiPOST(
		installationID,
		`
		mutation($input: CreateCheckRunInput!) {
			createCheckRun(input: $input) {
				clientMutationId
				checkRun {
					id
				}
			}
		}
		`,
		map[string]interface{}{
			"input": map[string]interface{}{
				"repositoryId": repositoryNodeID,
				"headSha":      headSHA,
				"name":         name,
				"status":       "COMPLETE",
				"conclusion":   "ACTION_REQUIRED",
				"output": map[string]string{
					"title":   "Summary",
					"summary": "Invalid Build",
					"text":    text,
				},
			},
		},
	)
	if D.IsErr(err) {
		log.Println("ActionRequiredCheckRun Error", err)
	}
}

// CreateCheckRun create CheckRun for build, return CheckRun ID
func CreateCheckRun(installationID, repositoryNodeID, headSHA, name string) (string, error) {
	if name == "" {
		name = "Build"
	}

	data, err := apiPOST(
		installationID,
		`
		mutation($input: CreateCheckRunInput!) {
			createCheckRun(input: $input) {
				clientMutationId
				checkRun {
					id
				}
			}
		}
		`,
		map[string]interface{}{
			"input": map[string]interface{}{
				"repositoryId": repositoryNodeID,
				"headSha":      headSHA,
				"name":         name,
				"status":       "QUEUED",
			},
		},
	)
	if D.IsErr(err) {
		return "", err
	}

	return data.CreateCheckRun.CheckRun.ID, nil
}

// QueuedCheckRun update CheckRun status to QUEUED
func QueuedCheckRun(installationID, repositoryNodeID, runID, buildID, logURL string) error {
	_, err := apiPOST(
		installationID,
		`
		mutation($input: UpdateCheckRunInput!) {
			updateCheckRun(input: $input) {
				clientMutationId
				checkRun {
					id
					repository {
						id
					}
				}
			}
		}
		`,
		map[string]interface{}{
			"input": map[string]interface{}{
				"checkRunId":   runID,
				"repositoryId": repositoryNodeID,
				"status":       "QUEUED",
				"detailsUrl":   logURL,
				"externalId":   buildID,
			},
		},
	)
	return err
}

// InProgressCheckRun update CheckRun status to IN_PROGESS
func InProgressCheckRun(installationID, repositoryNodeID, runID string) error {
	_, err := apiPOST(
		installationID,
		`
		mutation($input: UpdateCheckRunInput!) {
			updateCheckRun(input: $input) {
				clientMutationId
				checkRun {
					id
					repository {
						id
					}
				}
			}
		}
		`,
		map[string]interface{}{
			"input": map[string]interface{}{
				"checkRunId":   runID,
				"repositoryId": repositoryNodeID,
				"status":       "IN_PROGRESS",
			},
		},
	)
	return err
}

// CompletedCheckRun update CheckRun status to COMPLETED
func CompletedCheckRun(installationID, repositoryNodeID, runID string, conclusion CheckConclusionState, completedAt, summary, text string) error {
	_, err := apiPOST(
		installationID,
		`
		mutation($input: UpdateCheckRunInput!) {
			updateCheckRun(input: $input) {
				clientMutationId
				checkRun {
					id
					repository {
						id
					}
				}
			}
		}
		`,
		map[string]interface{}{
			"input": map[string]interface{}{
				"checkRunId":   runID,
				"repositoryId": repositoryNodeID,
				"status":       "COMPLETED",
				"completedAt":  completedAt,
				"conclusion":   conclusion,
				"output": map[string]string{
					"title":   "Summary",
					"summary": summary,
					"text":    text,
				},
			},
		},
	)
	return err
}
