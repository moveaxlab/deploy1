package argo

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/moveaxlab/deploy1/config"
	"github.com/moveaxlab/deploy1/output"
	log "github.com/sirupsen/logrus"
)

func Reset(service config.ServiceName, env config.Environment, imageTagParameter string, shouldWait bool) error {
	param := defaultImageTagParameter
	if imageTagParameter != "" {
		param = parameterName(imageTagParameter)
	}

	log.Infof("removing image tag override for service %s...", service)

	extraParams := config.Config.Argo.Environments[env].ArgoExtraParams

	cmdParams := append([]string{
		"--grpc-web",
		"app",
		"unset",
		string(service),
		"-p", string(param),
	},
		extraParams...)

	cmd := exec.Command(
		"argocd",
		cmdParams...,
	)
	customEnv := []string{
		fmt.Sprintf("ARGOCD_AUTH_TOKEN=%s", os.Getenv(config.Config.Argo.Environments[env].AuthTokenEnvVariable)),
		fmt.Sprintf("ARGOCD_SERVER=%s", config.Config.Argo.Environments[env].ServerName),
	}
	cmd.Env = append(os.Environ(), customEnv...)
	cmd.Stdout = output.OutLogger{}
	cmd.Stderr = output.ErrLogger{}
	log.Debugf("running command %s", cmd.String())

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to reset image tag for service %s: %w", service, err)
	}

	log.Infof("image tag override for service %s removed", service)

	if shouldWait {
		log.Infof("waiting for service %s to complete deployment...", service)
		err = wait(service, env)
		if err != nil {
			return err
		}
		log.Infof("service %s is now synced and healthy", service)
	}

	return nil
}
