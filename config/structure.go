package config

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
)

// a target environment for builds and deploys
type Environment string

// the name of a service
type Service string

// the name of an image
type ImageName string

// the name of a deployment
type ServiceName string

type BaseRegistryConfiguration struct {
	BasePath  string `json:"base_path"`
	Directory string `json:"directory"`
}

type RegistryConfiguration struct {
	BaseRegistryConfiguration
	Environments map[Environment]BaseRegistryConfiguration `json:"environments"`
}

type DockerConfiguration struct {
	// docker file to use for each build (relative to the project root)
	Dockerfile string `json:"dockerfile"`
}

type ScriptsConfiguration struct {
	// prepare script ran before all bundle + docker builds
	Prepare string `json:"prepare_bundle"`
	// bundle script ran before each docker build
	// receives in input the name of the service
	PreBuild string `json:"bundle"`
	// cleanup script ran after all images have been built
	Cleanup string `json:"post_bundle"`
}

type ServiceConfiguration struct {
	Directory         string               `json:"directory"`
	ServiceName       ServiceName          `json:"service_name"`
	ImageName         ImageName            `json:"image_name"`
	Scripts           ScriptsConfiguration `json:"scripts"`
	Dockerfile        string               `json:"dockerfile"`
	ImageTagParameter string               `json:"image_tag_parameter"`
}

type ArgoEnvironmentConfiguration struct {
	AuthTokenEnvVariable string   `json:"auth_token"`
	ServerName           string   `json:"server"`
	ServicePrefix        string   `json:"service_prefix"`
	ArgoExtraParams      []string `json:"argocli_extra_params"`
}

type ArgoConfiguration struct {
	Retries      int                                          `json:"retries"`
	Environments map[Environment]ArgoEnvironmentConfiguration `json:"environments"`
}

type Configuration struct {
	DefaultEnvironment Environment                      `json:"default_environment"`
	Argo               ArgoConfiguration                `json:"argo"`
	Docker             DockerConfiguration              `json:"docker"`
	Registry           RegistryConfiguration            `json:"registry"`
	Scripts            ScriptsConfiguration             `json:"scripts"`
	Services           map[Service]ServiceConfiguration `json:"services"`
}

var Config Configuration

func init() {
	log.SetFormatter(&log.TextFormatter{
		PadLevelText: true,
	})
	contents, err := ioutil.ReadFile("./deploy1.json")
	if err != nil {
		log.Debugf("no deploy1.json file found: %v", err)
		return
	}
	err = json.Unmarshal(contents, &Config)
	if err != nil {
		log.Debugf("failed to parse deploy1.json: %v", err)
		return
	}
}
