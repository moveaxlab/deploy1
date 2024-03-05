package bundle

import (
	"fmt"
	"os/exec"

	"github.com/moveaxlab/deploy1/config"
	"github.com/moveaxlab/deploy1/output"
	log "github.com/sirupsen/logrus"
)

func Bundle(service config.Service) error {
	bundleCmd, dir := config.GetPreBuildScript(service)

	if bundleCmd != "" {
		log.Infof("running bundle script for service %s...", service)
		cmd := exec.Command(bundleCmd, string(service))
		if dir != "" {
			cmd.Dir = dir
		}
		cmd.Stdout = output.OutLogger{}
		cmd.Stderr = output.OutLogger{}
		log.Debugf("running command %s", cmd.String())
		err := cmd.Run()
		if err != nil {
			return fmt.Errorf("bundle failed: %w", err)
		}
		log.Infof("service %s bundled", service)
	}

	return nil
}
