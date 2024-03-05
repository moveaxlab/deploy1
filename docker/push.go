package docker

import (
	"fmt"
	"github.com/moveaxlab/deploy1/config"
	"github.com/moveaxlab/deploy1/output"
	log "github.com/sirupsen/logrus"
	"os/exec"
)

func Push(service config.Service, env config.Environment, tag string) error {
	imageTag := config.GetImageTag(service, tag, env)
	log.Debugf("pushing image %s", imageTag)

	cmd := exec.Command(
		"docker",
		"push",
		imageTag,
	)

	log.Debugf("running command %s", cmd.String())
	cmd.Dir = config.Config.Services[service].Directory
	cmd.Stdout = output.OutLogger{}
	cmd.Stderr = output.ErrLogger{}

	err := cmd.Run()

	if err != nil {
		return fmt.Errorf("docker push failed: %w", err)
	}

	return nil
}
