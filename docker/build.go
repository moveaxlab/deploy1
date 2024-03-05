package docker

import (
	"fmt"
	"github.com/moveaxlab/deploy1/config"
	"github.com/moveaxlab/deploy1/output"
	log "github.com/sirupsen/logrus"
	"os/exec"
	"strings"
)

func Build(service config.Service, env config.Environment, tag string, buildArgs string, noCache bool) error {
	imageTag := config.GetImageTag(service, tag, env)
	log.Debugf("building image with tag %s", imageTag)

	dockerFile, err := config.GetDockerFile(service)
	if err != nil {
		return fmt.Errorf("unable to get dockerfile: %w", err)
	}
	log.Debugf("will use dockerfile %s", dockerFile)

	cmd := exec.Command(
		"docker",
		"build",
		"-t", imageTag,
		"-f", dockerFile,
		".",
	)

	actualBuildArgs := strings.Split(buildArgs, " ")

	for _, buildArg := range actualBuildArgs {
		if buildArg != "" {
			log.Debugf("adding custom build arg: %s", buildArg)
			cmd.Args = append(cmd.Args, "--build-arg", buildArg)
		}
	}

	if noCache {
		log.Debugf("disabling docker cache")
		cmd.Args = append(cmd.Args, "--no-cache")
	}

	log.Debugf("running command %s", cmd.String())
	cmd.Dir = config.Config.Services[service].Directory
	cmd.Stdout = output.OutLogger{}
	cmd.Stderr = output.ErrLogger{}

	err = cmd.Run()

	if err != nil {
		return fmt.Errorf("docker build failed: %w", err)
	}

	return nil
}
