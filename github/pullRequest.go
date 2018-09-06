package github

import (
	"strings"

	D "github.com/NeoJRotary/describe-go"
)

// PullRequest pull request
type PullRequest struct {
	ID          string `json:"id"`
	Number      int    `json:"number"`
	BaseRefName string `json:"baseRefName"`
}

// GetPullRequests get pull request list
func GetPullRequests(installationID, repoFullName string) ([]PullRequest, error) {
	ss := strings.Split(repoFullName, "/")

	data, err := apiPOST(
		installationID,
		`
		{
			repository(owner:"`+ss[0]+`", name:"`+ss[1]+`") {
				pullRequests(first:100, states:OPEN) {
					nodes{
						id
						number
						baseRefName
					}
				}
			}
		}
		`,
		nil,
	)
	if D.IsErr(err) {
		return nil, err
	}

	return data.Repository.PullRequests.Nodes, nil
}

// GetAssociatedPullRequests get associated pull request list
func GetAssociatedPullRequests(installationID, repoFullName, headBranch string) ([]PullRequest, error) {
	ss := strings.Split(repoFullName, "/")

	data, err := apiPOST(
		installationID,
		`
		{
			repository(owner: "`+ss[0]+`", name: "`+ss[1]+`") {
				ref(qualifiedName: "refs/heads/`+headBranch+`") {
					associatedPullRequests(first: 100, states: OPEN) {
						nodes {
							id
							number
							baseRefName
						}
					}
				}
			}
		}
		`,
		nil,
	)
	if D.IsErr(err) {
		return nil, err
	}

	return data.Repository.Ref.AssociatedPullRequests.Nodes, nil
}
