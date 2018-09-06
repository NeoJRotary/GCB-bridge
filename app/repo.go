package app

import (
	"io/ioutil"
	"log"
	"path"

	D "github.com/NeoJRotary/describe-go"
	exec "github.com/NeoJRotary/exec-go"
)

// InitRepo init repo
func InitRepo() {
	// prepare git folder
	_, err := exec.RunCmd("/", "mkdir", "-p", GitFolder)
	if D.IsErr(err) {
		log.Fatal(err)
	}

}

// Repo repo struct
type Repo struct {
	Event            string
	InstallationID   string
	RepositoryNodeID string
	Dir              string
	FullName         string
	Branch           string
	AssociatedBases  []string
	BaseBranch       string
	Tag              string
	BeforeSHA        string
	AfterSHA         string
	BuildFilePath    string
	Removed          bool
}

// LogError repo log error
func (repo *Repo) LogError(title string, err error) bool {
	text := `----- app.Repo LogError -----
InstallationID: ` + repo.InstallationID + `
RepositoryNodeID: ` + repo.RepositoryNodeID + `
FullName: ` + repo.FullName + `
Branch: ` + repo.Branch + `
Tag: ` + repo.Tag + `
CommitSHA: ` + repo.BeforeSHA + ` >>> ` + repo.AfterSHA + `
`

	log.Println(text, title, err)
	return false
}

// Init init repo
func (repo *Repo) Init() bool {
	if repo.Branch == "" && repo.Tag == "" {
		return repo.LogError("Init", D.NewErr("need one of branch/tag"))
	}

	cloneTarget := ""
	if repo.Branch != "" {
		cloneTarget = repo.Branch
	}
	if repo.Tag != "" {
		cloneTarget = repo.Tag
	}

	tempDir, err := ioutil.TempDir(GitFolder, "")
	if D.IsErr(err) {
		return repo.LogError("Init", err)
	}

	_, err = exec.RunCmd(tempDir, "git", "clone", "-b", cloneTarget, "https://x-access-token:"+GetAccessToken(repo.InstallationID)+"@github.com/"+repo.FullName+".git", ".")
	if D.IsErr(err) {
		return repo.LogError("Init", err)
	}

	buildFile := path.Join(tempDir, CloudBuildFileName)
	if !D.FileExist(buildFile) {
		repo.Remove()
		return repo.LogError("Init", D.NewErr(CloudBuildFileName+" Not found."))
	}

	// setup
	repo.Dir = tempDir
	repo.BuildFilePath = buildFile

	if repo.Branch != "" {
		// get initial commit as beforeSHA
		if repo.BeforeSHA == "" {
			out, err := exec.RunCmd(repo.Dir, "git", "rev-list", "--max-parents=0", "HEAD")
			if D.IsErr(err) {
				return repo.LogError("Init", err)
			}
			repo.BeforeSHA = D.String(out).Trim("\n").TrimSpace().Get()
		}

		// get HEAD as afterSHA
		if repo.AfterSHA == "" {
			out, err := exec.RunCmd(repo.Dir, "git", "rev-parse", "HEAD")
			if D.IsErr(err) {
				return repo.LogError("Init", err)
			}
			repo.AfterSHA = D.String(out).Trim("\n").TrimSpace().Get()
		}
	}

	if repo.Tag != "" {
		// setup both BeforeSHA and AfterSHA to tag commit
		out, err := exec.RunCmd(repo.Dir, "git", "rev-parse", "HEAD")
		if D.IsErr(err) {
			return repo.LogError("Init", err)
		}
		sha := D.String(out).Trim("\n").TrimSpace().Get()
		repo.BeforeSHA = sha
		repo.AfterSHA = sha
	}

	return true
}

// GetChanges get changes filename between BeforeSHA and AfterSHA
func (repo *Repo) GetChanges() []string {
	out, err := exec.RunCmd(repo.Dir, "git", "diff", "--name-only", repo.BeforeSHA, repo.AfterSHA)
	if D.IsErr(err) {
		repo.LogError("GetChanges", err)
		return nil
	}

	list := D.String(out).Trim("\n").Split("\n").ElmTrimSpace().Get()
	for i, f := range list {
		list[i] = path.Join(repo.Dir, f)
	}
	return list
}

// Remove remove repo
func (repo *Repo) Remove() {
	if repo.Removed {
		return
	}
	repo.Removed = true

	exec.RunCmd("/", "rm", "-rf", repo.Dir)
	repo.Dir = ""
}
