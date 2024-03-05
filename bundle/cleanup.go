package bundle

import (
	"fmt"
	"github.com/moveaxlab/deploy1/config"
	"github.com/moveaxlab/deploy1/output"
	log "github.com/sirupsen/logrus"
	"os/exec"
)

func Cleanup() error {
	cleanupCmd := config.Config.Scripts.Cleanup

	if cleanupCmd != "" {
		log.Infof("running post-bundle cleanup script...")
		cmd := exec.Command(cleanupCmd)
		cmd.Stdout = output.OutLogger{}
		cmd.Stderr = output.OutLogger{}
		log.Debugf("running command %s", cmd.String())
		err := cmd.Run()
		if err != nil {
			return fmt.Errorf("cleanup failed: %w", err)
		}
		log.Infof("post-bundle cleanup complete")
	}

	return nil
}
