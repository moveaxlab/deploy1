package config

import (
	"fmt"
	"github.com/spf13/cobra"
	"strings"
)

func GetServiceName(service Service, environment Environment) ServiceName {
	servicePrefix := Config.Argo.Environments[environment].ServicePrefix
	if servicePrefix != "" {
		return ServiceName(fmt.Sprintf("%s-%s", servicePrefix, Config.Services[service].ServiceName))
	} else {
		return Config.Services[service].ServiceName
	}
}

func GetImageName(service Service) ImageName {
	return Config.Services[service].ImageName
}

func GetImageTagParameter(service Service) string {
	return Config.Services[service].ImageTagParameter
}

func GetAllServices() []Service {
	res := make([]Service, 0, len(Config.Services))

	for e := range Config.Services {
		res = append(res, e)
	}

	return res
}

func stringSliceContains(slice []string, elem string) bool {
	for _, v := range slice {
		if v == elem {
			return true
		}
	}
	return false
}

func AutocompleteService(alreadyCompleted []string, toComplete string) []string {
	cobra.CompDebugln(fmt.Sprintf("autocompleting service \"%s\"", toComplete), true)
	res := make([]string, 0, len(Config.Services))

	for e := range Config.Services {
		if stringSliceContains(alreadyCompleted, string(e)) {
			cobra.CompDebugln(fmt.Sprintf("service \"%s\" already chosen", string(e)), true)
			continue
		}
		if strings.HasPrefix(string(e), toComplete) {
			cobra.CompDebugln(fmt.Sprintf("service \"%s\" matched", string(e)), true)
			res = append(res, string(e))
		}
	}

	return res
}
