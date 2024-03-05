package config

func GetPreBuildScript(service Service) (string, string) {
	if Config.Services[service].Scripts.PreBuild != "" {
		return Config.Services[service].Scripts.PreBuild, Config.Services[service].Directory
	}

	return Config.Scripts.PreBuild, ""
}
