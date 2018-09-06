package trigger

import (
	"sync"

	"github.com/NeoJRotary/GCB-bridge/app"
	"github.com/NeoJRotary/GCB-bridge/gcloud"
)

// EventHandler trigger event handler
func EventHandler(repo *app.Repo) {
	if !repo.Init() {
		return
	}

	defer repo.Remove()

	builds := getValidBuilds(repo)
	if len(builds) != 0 {
		var wg sync.WaitGroup
		wg.Add(len(builds))
		for _, build := range builds {
			// start build if available
			if build.Available {
				go gcloud.StartBuild(&wg, repo, build.Name, build.File)
			}
		}
		wg.Wait()
	}
}
