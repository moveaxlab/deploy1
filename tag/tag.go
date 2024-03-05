package tag

import (
	"fmt"
	"github.com/moveaxlab/deploy1/git"
	"github.com/moveaxlab/deploy1/utils"
	log "github.com/sirupsen/logrus"
)

func GetTag(tag string) (string, error) {
	if tag != "" {
		log.Debugf("valid tag provided")
		return tag, nil
	}

	currentBranch, err := git.CurrentBranch()
	if err != nil {
		return "", fmt.Errorf("failed to create tag: %w", err)
	}
	log.Debugf("current branch is %s", currentBranch)

	if currentBranch == "dev" {
		return "dev", nil
	}

	taskName := utils.GetTaskName(currentBranch)

	if taskName == "" {
		return "", fmt.Errorf("current branch does not contain a valid tag, please provide manually")
	}

	return taskName, nil
}
