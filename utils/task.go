package utils

import "regexp"

func GetTaskName(branchName string) string {
	re := regexp.MustCompile("[A-Z]+-[0-9]+")
	return re.FindString(branchName)
}
