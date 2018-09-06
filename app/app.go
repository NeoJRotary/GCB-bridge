package app

// GitFolder default location to put temp git repos
var GitFolder = "/gcb-bridge/git/"

// CloudBuildFileName default yaml file name for building
var CloudBuildFileName = "cloudbuild.bridge.yaml"

// Init init app pacakge
func Init() {
	InitRepo()
	InitToken()
}
