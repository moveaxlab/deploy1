package config

import "fmt"

func ValidateServiceName(args []string) ([]Service, error) {
	res := make([]Service, 0, len(args))

	if len(args) == 0 {
		return nil, fmt.Errorf("please provide at least one service name")
	}

	for _, serviceName := range args {
		if _, ok := Config.Services[Service(serviceName)]; !ok {
			return nil, fmt.Errorf("unknown service '%s'", serviceName)
		}
		res = append(res, Service(serviceName))
	}

	return res, nil
}

func ValidateEnvironment(env string) (Environment, error) {
	if _, ok := Config.Registry.Environments[Environment(env)]; !ok {
		return "", fmt.Errorf("unknown environment '%s'", env)
	}

	return Environment(env), nil
}
