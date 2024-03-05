package config

import (
	"fmt"
	"path/filepath"
)

func GetDockerFile(service Service) (string, error) {
	if Config.Services[service].Dockerfile != "" {
		return Config.Services[service].Dockerfile, nil
	}

	dockerFile := Config.Docker.Dockerfile
	dockerFile, err := filepath.Abs(dockerFile)
	if err != nil {
		return "", fmt.Errorf("unable to resolve root dockerfile: %w", err)
	}

	return dockerFile, nil
}
