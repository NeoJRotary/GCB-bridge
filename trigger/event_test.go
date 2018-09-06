package trigger

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/NeoJRotary/GCB-bridge/app"
	D "github.com/NeoJRotary/describe-go"
	"gopkg.in/yaml.v2"
)

// test push to branch
func TestTrigger_Branch(t *testing.T) {
	initTest()

	repo := &app.Repo{
		Event:          "Branch",
		InstallationID: "",
		FullName:       "NeoJRotary/GCB-bridge-test",
		Branch:         "feature/newAPI",
		// commit `only readme`
		BeforeSHA: "0030064bce54a1051d4642771caed3b8edf2ad21",
		// commit `update test yaml`
		AfterSHA: "a96e6eb1a0a51508b7de8a23444f40e00aaf31e3",
	}

	triggerTester(t, repo)
}

// test push a tag
func TestTrigger_Tag(t *testing.T) {
	initTest()

	repo := &app.Repo{
		Event:          "Tag",
		InstallationID: "",
		FullName:       "NeoJRotary/GCB-bridge-test",
		Tag:            "v1.3.1",
	}

	triggerTester(t, repo)
}

// test pull request opened
func TestTrigger_PullRequest(t *testing.T) {
	initTest()

	repo := &app.Repo{
		Event:          "PullRequest",
		InstallationID: "",
		FullName:       "NeoJRotary/GCB-bridge-test",
		Branch:         "develop",
		BaseBranch:     "master",
	}

	triggerTester(t, repo)
}

// test push to branch which is a head of pull request and there is build be triggered by it's base
func TestTrigger_AssociatedBase(t *testing.T) {
	initTest()

	repo := &app.Repo{
		Event:           "Branch",
		InstallationID:  "",
		FullName:        "NeoJRotary/GCB-bridge-test",
		Branch:          "develop",
		AssociatedBases: []string{"master"},
	}

	triggerTester(t, repo)
}

func initTest() {
	DEBUG = true
	app.InitRepo()
}

func triggerTester(t *testing.T, repo *app.Repo) {
	if !repo.Init() {
		t.Fatal("repo init failed")
	}

	builds := getValidBuilds(repo)
	if len(builds) != 1 {
		t.Fatal("should get one")
	}

	b, err := ioutil.ReadFile(builds[0].File)
	if D.IsErr(err) {
		t.Fatal(err)
	}

	// print the file itself for checking
	fmt.Println("----------------------")
	fmt.Println(string(b))
	fmt.Println("----------------------")
	// should receive:
	//
	// artifacts:
	// 	objects:
	// 		location: 'gs://$PROJECT_ID/'
	// 		paths: ['hello']
	// steps:
	// - name: 'gcr.io/cloud-builders/go'
	// 	args: ['test', '-run', 'Success']
	// - name: 'gcr.io/cloud-builders/go'
	// 	args: ['test', '-run', 'Fail']
	// 	env: ['PROJECT_ROOT=hello']

	var obj map[interface{}]interface{}
	err = yaml.Unmarshal(b, &obj)
	if D.IsErr(err) {
		t.Fatal(err)
	}

	if len(obj) != 3 {
		t.Fatal("top level should have 3 key")
	}

	// check options
	_, ok := obj["options"]
	if !ok {
		t.Fatal("top level should have key: options")
	}
	options, ok := obj["options"].(map[interface{}]interface{})
	if !ok {
		t.Fatal("options should be map[]")
	}
	_, ok = options["substitution_option"]
	if !ok {
		t.Fatal("options should have substitution_option")
	}

	// check artifacts
	_, ok = obj["artifacts"]
	if !ok {
		t.Fatal("top level should have key: artifacts")
	}
	artifacts, ok := obj["artifacts"].(map[interface{}]interface{})
	if !ok {
		t.Fatal("artifacts should be map[]")
	}
	_, ok = artifacts["objects"]
	if !ok {
		t.Fatal("artifacts should have objects")
	}
	objects, ok := artifacts["objects"].(map[interface{}]interface{})
	if !ok {
		t.Fatal("artifacts.objects should be map[]")
	}
	_, ok = objects["location"]
	if !ok {
		t.Fatal("artifacts.objects should have location")
	}
	_, ok = objects["location"].(string)
	if !ok {
		t.Fatal("artifacts.objects.location should be string")
	}
	_, ok = objects["paths"]
	if !ok {
		t.Fatal("artifacts.objects should have paths")
	}
	_, ok = objects["paths"].([]interface{})
	if !ok {
		t.Fatal("artifacts.objects.paths should be array")
	}

	// check steps
	_, ok = obj["steps"]
	if !ok {
		t.Fatal("top level should have key: steps")
	}
	steps, ok := obj["steps"].([]interface{})
	if !ok {
		t.Fatal("steps should be array")
	}
	for i, elm := range steps {
		step, ok := elm.(map[interface{}]interface{})
		if !ok {
			t.Fatal("steps elements should be map[]")
		}

		_, ok = step["name"].(string)
		if !ok {
			t.Fatal("step.name should be string")
		}
		_, ok = step["args"].([]interface{})
		if !ok {
			t.Fatal("step.args should be array")
		}

		if i == 1 {
			_, ok = step["env"].([]interface{})
			if !ok {
				t.Fatal("step[1].env should be array")
			}
		}
	}
}
