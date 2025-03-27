package argo

import (
	"encoding/json"
	"fmt"
	"github.com/moveaxlab/deploy1/config"
	"github.com/moveaxlab/deploy1/output"
	log "github.com/sirupsen/logrus"
	"os"
	"os/exec"
)

type parameterName string

type resourceGroup string

type resourceKind string

const (
	defaultImageTagParameter parameterName = "image.tag"
	defaultImageTag                        = "latest"
	deploymentResource       resourceKind  = "Deployment"
	statefulSetResource      resourceKind  = "StatefulSet"
	appResourceGroup         resourceGroup = "apps"
)

type argoParameters struct {
	Name  parameterName `json:"name"`
	Value string        `json:"value"`
}

type argoParamsHelm struct {
	Parameters []argoParameters `json:"parameters"`
}

type argoParamsSource struct {
	Helm argoParamsHelm `json:"helm"`
}

type argoParamsSpec struct {
	Source argoParamsSource `json:"source"`
}

type argoParamsResource struct {
	Group resourceGroup `json:"group"`
	Kind  resourceKind  `json:"kind"`
}

type argoParamsStatus struct {
	Resources []argoParamsResource `json:"resources"`
}

type argoParams struct {
	Spec   argoParamsSpec   `json:"spec"`
	Status argoParamsStatus `json:"status"`
}

type serviceInfo struct {
	currentImageTag string
	resourceKind    resourceKind
}

func getServiceInfo(service config.ServiceName, env config.Environment, customImageTagParameter string) (*serviceInfo, error) {
	cmd := exec.Command(
		"argocd",
		"--grpc-web",
		"app",
		"get",
		string(service),
		"show-params",
		"-o", "json",
		"--insecure",
		"--plaintext",
		"--loglevel=debug",
	)
	cmd.Env = []string{
		fmt.Sprintf("ARGOCD_AUTH_TOKEN=%s", os.Getenv(config.Config.Argo.Environments[env].AuthTokenEnvVariable)),
		fmt.Sprintf("ARGOCD_SERVER=%s", config.Config.Argo.Environments[env].ServerName),
	}

	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	log.Debugf("running command %s", cmd.String())

	res, err := cmd.Output()

	log.Debugf("output:\n%s", string(res))
	if err != nil {
		return nil, fmt.Errorf("failed to get current tag: %w", err)
	}

	var response argoParams

	err = json.Unmarshal(res, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to parse argocd response: %w", err)
	}

	parameters := response.Spec.Source.Helm.Parameters

	imageTag := defaultImageTag
	imageTagParameter := defaultImageTagParameter

	if customImageTagParameter != "" {
		imageTagParameter = parameterName(customImageTagParameter)
	}

	for _, param := range parameters {
		if param.Name == imageTagParameter {
			imageTag = param.Value
		}
	}

	var resourceKind resourceKind

	for _, resource := range response.Status.Resources {
		if resource.Group == appResourceGroup {
			resourceKind = resource.Kind
			break
		}
	}

	if resourceKind == "" {
		return nil, fmt.Errorf("unable to determine resource kind for service %s", service)
	}

	log.Debugf("no image tag override found on argo")
	return &serviceInfo{
		currentImageTag: imageTag,
		resourceKind:    resourceKind,
	}, nil
}

func restart(service config.ServiceName, env config.Environment, kind resourceKind) error {
	cmd := exec.Command(
		"argocd",
		"--grpc-web",
		"app",
		"actions",
		"run",
		string(service),
		"restart",
		"--kind", string(kind),
		"--all",
		"--insecure",
		"--plaintext",
	)
	cmd.Env = []string{
		fmt.Sprintf("ARGOCD_AUTH_TOKEN=%s", os.Getenv(config.Config.Argo.Environments[env].AuthTokenEnvVariable)),
		fmt.Sprintf("ARGOCD_SERVER=%s", config.Config.Argo.Environments[env].ServerName),
	}
	cmd.Stdout = output.OutLogger{}
	cmd.Stderr = output.ErrLogger{}
	log.Debugf("running command %s", cmd.String())

	err := cmd.Run()

	if err != nil {
		return fmt.Errorf("failed to restart service %s: %w", service, err)
	}

	return nil
}

func deploy(service config.ServiceName, tag string, env config.Environment, customImageTagParameter string) error {
	imageTagParameter := defaultImageTagParameter
	if customImageTagParameter != "" {
		imageTagParameter = parameterName(customImageTagParameter)
	}
	cmd := exec.Command(
		"argocd",
		"--grpc-web",
		"app",
		"set",
		string(service),
		"--helm-set-string", fmt.Sprintf("%s=%s", imageTagParameter, tag),
		"--insecure",
		"--plaintext",
	)
	cmd.Env = []string{
		fmt.Sprintf("ARGOCD_AUTH_TOKEN=%s", os.Getenv(config.Config.Argo.Environments[env].AuthTokenEnvVariable)),
		fmt.Sprintf("ARGOCD_SERVER=%s", config.Config.Argo.Environments[env].ServerName),
	}
	cmd.Stdout = output.OutLogger{}
	cmd.Stderr = output.ErrLogger{}
	log.Debugf("running command %s", cmd.String())

	err := cmd.Run()

	if err != nil {
		return fmt.Errorf("failed to restart service %s: %w", service, err)
	}

	return nil
}

func Deploy(service config.ServiceName, tag string, env config.Environment, imageTagParameter string) error {
	log.Infof("retrieving current tag for service %s from argo...", service)
	serviceInfo, err := getServiceInfo(service, env, imageTagParameter)
	if err != nil {
		return fmt.Errorf("failed to get current tag of service %s: %w", service, err)
	}
	log.Infof("current tag for service %s is %s", service, serviceInfo.currentImageTag)
	log.Infof("service %s is a %s", service, serviceInfo.resourceKind)

	if serviceInfo.currentImageTag == tag {
		log.Infof("current tag and given tag are identical, restarting...")
		err = restart(service, env, serviceInfo.resourceKind)
		if err != nil {
			return fmt.Errorf("restart of service %s failed: %w", service, err)
		}
		log.Infof("service %s restarted with tag %s", service, tag)
	} else {
		log.Infof("current tag and given tag differ, overriding...")
		err = deploy(service, tag, env, imageTagParameter)
		if err != nil {
			return fmt.Errorf("deploy of service %s failed: %w", service, err)
		}
		log.Infof("override of service %s to tag %s complete", service, tag)
	}

	return nil
}
