package config

import "strings"

func AutocompleteEnvironment(toComplete string) []string {
	res := make([]string, 0, len(Config.Registry.Environments))

	for e := range Config.Registry.Environments {
		if strings.HasPrefix(string(e), toComplete) {
			res = append(res, string(e))
		}
	}

	return res
}
