package docker

import (
	"bytes"
	"fmt"
	"github.com/moveaxlab/deploy1/config"
	log "github.com/sirupsen/logrus"
	"os/exec"
	"strings"
)

func TagExists(service config.Service, env config.Environment, tag string) (bool, error) {
	imageTag := config.GetImageTag(service, tag, env)
	log.Debugf("checking if %s exists...", imageTag)

	stderr := bytes.NewBuffer([]byte{})

	// https://docs.docker.com/engine/reference/commandline/manifest/
	cmd := exec.Command(
		"docker",
		"manifest",
		"inspect",
		imageTag,
	)
	cmd.Env = []string{
		"DOCKER_CLI_EXPERIMENTAL=enabled",
	}
	cmd.Stderr = stderr

	log.Debugf("running command %s", cmd.String())
	output, err := cmd.Output()
	log.Debugf("output: %s", string(output))
	log.Debugf("error: %s", cmd.Stderr)

	if err != nil {
		if strings.HasPrefix(stderr.String(), "no such manifest") {
			return false, nil
		} else {
			return false, fmt.Errorf("failed to check if image exists: %w", err)
		}
	}

	return true, nil
}
