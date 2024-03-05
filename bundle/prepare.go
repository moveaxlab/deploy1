package bundle

import (
	"fmt"
	"os/exec"

	"github.com/moveaxlab/deploy1/config"
	"github.com/moveaxlab/deploy1/output"
	log "github.com/sirupsen/logrus"
)

func Prepare() error {
	prepareCmd := config.Config.Scripts.Prepare

	if prepareCmd != "" {
		log.Infof("running prepare bundle script...")
		cmd := exec.Command(prepareCmd)
		cmd.Stdout = output.OutLogger{}
		cmd.Stderr = output.OutLogger{}
		log.Debugf("running command %s", cmd.String())
		err := cmd.Run()
		if err != nil {
			return fmt.Errorf("prepare failed: %w", err)
		}
		log.Infof("prepare finished")
	}

	return nil
}
