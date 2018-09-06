package gcloud

import (
	"io/ioutil"
	"log"
	"os"
	"sync"

	// "github.com/NeoJRotary/GCB-bridge/git"
	"github.com/NeoJRotary/GCB-bridge/app"
	"github.com/NeoJRotary/GCB-bridge/github"
	D "github.com/NeoJRotary/describe-go"
	exec "github.com/NeoJRotary/exec-go"
)

var (
	// ProjectID GCP project ID
	ProjectID string
	// StorageBucket GCS Bucket for source and log
	StorageBucket string
)

// SkipGithubAPI skip github API for glcoud unit test
var SkipGithubAPI = false

// Init init for gcloud
func Init() {
	initCredentials()
	initPubSub()
}

func initCredentials() {
	SA := D.GetENV("GCLOUD_SERVICE_ACCOUNT", "")
	if SA == "" {
		log.Fatal("GCLOUD_SERVICE_ACCOUNT should not be empty")
	}

	ProjectID = D.GetENV("GCLOUD_PROJECT_ID", "")
	if ProjectID == "" {
		log.Fatal("GCLOUD_PROJECT_ID should not be empty")
	}

	StorageBucket = D.GetENV("GCLOUD_STORAGE_BUCKET", "")
	if StorageBucket == "" {
		log.Fatal("GCLOUD_STORAGE_BUCKET should not be empty")
	}

	f, err := ioutil.TempFile("", "credentials-*.json")
	if D.IsErr(err) {
		log.Fatal("Create service account tempFile failed", err)
	}
	_, err = f.Write([]byte(SA))
	if D.IsErr(err) {
		log.Fatal("Create service account tempFile failed", err)
	}
	f.Close()
	os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", f.Name())

	_, err = exec.RunCmd("", "gcloud", "auth", "activate-service-account", "--key-file", f.Name())
	if D.IsErr(err) {
		log.Fatal(err)
	}
	_, err = exec.RunCmd("", "gcloud", "config", "set", "project", ProjectID)
	if D.IsErr(err) {
		log.Fatal(err)
	}

	log.Println("GCloud SDK Inited")
}

// StartBuild start build
func StartBuild(wg *sync.WaitGroup, repo *app.Repo, name, buildFile string) {
	defer wg.Done()

	log.Println("StartBuild", repo.FullName, repo.Branch, repo.Tag)

	var (
		runID string
		err   error
	)
	if !SkipGithubAPI {
		runID, err = github.CreateCheckRun(repo.InstallationID, repo.RepositoryNodeID, repo.AfterSHA, name)
		if D.IsErr(err) {
			log.Println("StartBuild error", err)
			return
		}
	}

	subsList := D.StringSlice()
	subsList.Push("REPO_NAME=" + repo.FullName)
	subsList.Push("BRANCH_NAME=" + repo.Branch)
	subsList.Push("TAG_NAME=" + repo.Tag)
	subsList.Push("COMMIT_SHA=" + repo.AfterSHA)
	subsList.Push("_GITHUB_INSTALLATION_ID=" + repo.InstallationID)
	subsList.Push("_GITHUB_REPOSITORY_NODE_ID=" + repo.RepositoryNodeID)
	subsList.Push("_GITHUB_CHECKRUN_ID=" + runID)

	if EnableMsgListener {
		subsList.Push("_BRIDGE_UID=" + MsgListener.uid)
	}

	out, err := exec.RunCmd(
		"",
		"gcloud", "builds", "submit", repo.Dir,
		"--config", buildFile,
		"--substitutions", subsList.Join(",").Get(),
		"--async",
		"--format", "value(id)",
		"--gcs-log-dir", "gs://"+StorageBucket+"/log",
		"--gcs-source-staging-dir", "gs://"+StorageBucket+"/source",
	)
	if D.IsErr(err) {
		log.Println("StartBuild failed", err)
		return
	}

	log.Println("StartBuild ID:", D.String(out).Trim("\n").Get())
}
