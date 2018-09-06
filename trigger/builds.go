package trigger

import (
	"bytes"
	"io/ioutil"
	"log"
	"strconv"

	"github.com/NeoJRotary/GCB-bridge/app"
	D "github.com/NeoJRotary/describe-go"
	yaml "gopkg.in/yaml.v2"
)

// cloudBuild one GCB cloudbuild.yaml configuration. (one steps set)
type cloudBuild struct {
	repo   *app.Repo
	b      []byte
	rawObj map[string]interface{}

	Log       *D.TypeStringSlice
	File      string
	Available bool

	Name     string           `yaml:"name"`
	Triggers []trigger        `yaml:"triggers"`
	Steps    []cloudBuildStep `yaml:"steps"`
}

type cloudBuildStep struct {
	Name     string    `yaml:"name"`
	Triggers []trigger `yaml:"triggers"`
	unused   bool
}

func getValidBuilds(repo *app.Repo) []*cloudBuild {
	b, err := ioutil.ReadFile(repo.BuildFilePath)
	if D.IsErr(err) {
		log.Println("getValidBuilds ReadFile Error " + err.Error())
		return nil
	}

	// split builds
	bb := bytes.Split(b, []byte("---\n"))

	builds := []*cloudBuild{}
	for i, b := range bb {
		cb := parseBuild(repo, b)

		title := "Builds[" + strconv.Itoa(i) + "] " + cb.Name
		if cb.Available {
			cb.Log.Shift(title + " AVAILABLE")
			builds = append(builds, cb)
		} else {
			cb.Log.Shift(title + " UNAVAILABLE")
		}

		log.Println(cb.Log.Join("\n").Get())
	}

	return builds
}

// parseBuild parse one build
func parseBuild(repo *app.Repo, b []byte) *cloudBuild {
	cb := &cloudBuild{
		repo: repo,
		b:    b,
		Log:  D.StringSlice(),
	}

	ok := cb.checkRawObj()
	if !ok {
		return cb
	}

	ok = cb.unmarshalConf()
	if !ok {
		return cb
	}

	ok = cb.validate()
	if !ok {
		return cb
	}

	ok = cb.saveFile()
	if !ok {
		return cb
	}

	cb.Available = true
	return cb
}

// checkRawObj get raw map[string]interface{} and check structure
func (cb *cloudBuild) checkRawObj() bool {
	var raw map[string]interface{}
	err := yaml.Unmarshal(cb.b, &raw)
	if D.IsErr(err) {
		cb.Log.Push("Raw unmarshal failed: " + err.Error())
		return false
	}

	temp, ok := raw["steps"]
	if !ok {
		cb.Log.Push("Invalid Build: should have [steps]")
		return false
	}

	steps, ok := temp.([]interface{})
	if !ok {
		cb.Log.Push("Invalid Build: [steps] should be array")
		return false
	}

	for _, step := range steps {
		_, ok = step.(map[interface{}]interface{})
		if !ok {
			cb.Log.Push("Invalid Build: [steps] step should be object")
			return false
		}
	}

	cb.rawObj = raw

	return true
}

// unmarshalConf unmarshal yaml to struct
func (cb *cloudBuild) unmarshalConf() bool {
	err := yaml.Unmarshal(cb.b, cb)
	if D.IsErr(err) {
		cb.Log.Push("Conf unmarshal failed: " + err.Error())
		return false
	}

	return true
}

// validate validate both global and step triggers
func (cb *cloudBuild) validate() bool {
	if !cb.validTriggers(cb.Triggers) {
		cb.Log.Push("Invlaid global triggers")
		return false
	}

	count := 0
	for i, step := range cb.Steps {
		if cb.validTriggers(step.Triggers) {
			count++
		} else {
			// ignore the step
			cb.rawObj["steps"].([]interface{})[i] = nil
		}
	}

	if count == 0 {
		cb.Log.Push("No any available steps after validation")
		return false
	}

	return true
}

// validTriggers valid one triggers set
func (cb *cloudBuild) validTriggers(triggers []trigger) bool {
	if len(triggers) == 0 {
		return true
	}

	// go through all triggers, if any of them is valid, return true
	for _, trg := range triggers {
		if trg.isValid(cb) {
			return true
		}
	}

	return false
}

// reformRaw remove keys that only be used by GCB-bridge
func (cb *cloudBuild) reformRaw() {
	delete(cb.rawObj, "triggers")
	delete(cb.rawObj, "name")

	newSteps := []map[interface{}]interface{}{}
	for _, step := range cb.rawObj["steps"].([]interface{}) {
		if step == nil {
			continue
		}
		stepObj := step.(map[interface{}]interface{})
		delete(stepObj, "triggers")
		newSteps = append(newSteps, stepObj)
	}
	cb.rawObj["steps"] = newSteps

	// add options for self-substitution skip error
	cb.rawObj["options"] = map[string]string{"substitution_option": "ALLOW_LOOSE"}
}

// saveFile save reformed file
func (cb *cloudBuild) saveFile() bool {
	cb.reformRaw()

	b, err := yaml.Marshal(cb.rawObj)
	if D.IsErr(err) {
		cb.Log.Push("Marshal RawObj Error " + err.Error())
		return false
	}

	f, err := ioutil.TempFile(cb.repo.Dir, "cloudbuild.*.yaml")
	if D.IsErr(err) {
		cb.Log.Push("ioutil.TempFile Error " + err.Error())
		return false
	}

	defer f.Close()

	_, err = f.Write(b)
	if D.IsErr(err) {
		cb.Log.Push("ioutil.TempFile.Write Error " + err.Error())
		return false
	}

	cb.File = f.Name()

	return true
}
