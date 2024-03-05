package config

import "fmt"

func GetImageTag(service Service, tag string, environment Environment) string {
	return fmt.Sprintf(
		"%s/%s/%s:%s",
		Config.Registry.BasePath,
		Config.Registry.Environments[environment].Directory,
		GetImageName(service),
		tag,
	)
}
