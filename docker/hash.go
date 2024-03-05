package docker

import (
	"fmt"
	"github.com/moveaxlab/deploy1/config"
	log "github.com/sirupsen/logrus"
	"os/exec"
	"strings"
)

func GetHash(service config.Service, env config.Environment, tag string) (string, error) {
	imageTag := config.GetImageTag(service, tag, env)
	log.Debugf("retrieving hash of image %s", imageTag)

	cmd := exec.Command(
		"docker",
		"inspect",
		"--format='{{index .RepoDigests 0}}'",
		imageTag,
	)

	res, err := cmd.Output()

	if err != nil {
		return "", fmt.Errorf("failed to get image hash for service %s: %w", service, err)
	}

	hash := strings.Split(strings.Trim(string(res), "'"), "@")[1]

	return hash, nil
}
