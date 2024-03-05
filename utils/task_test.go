package utils

import (
	"testing"
)

func TestGetTaskName(t *testing.T) {
	cases := []struct {
		branchName string
		expected   string
	}{
		{
			branchName: "feature/MX-1234-my-feature",
			expected:   "MX-1234",
		},
		{
			branchName: "MX-22-my-task",
			expected:   "MX-22",
		},
		{
			branchName: "CS-666-my-task-2",
			expected:   "CS-666",
		},
		{
			branchName: "bugfix/ISD-123-some-bug",
			expected:   "ISD-123",
		},
		{
			branchName: "bad-branch-name",
			expected:   "",
		},
		{
			branchName: "MX-1234/MX-6543-double-task",
			expected:   "MX-1234",
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.branchName, func(t *testing.T) {
			actual := GetTaskName(testCase.branchName)

			if actual != testCase.expected {
				t.Errorf("GetTaskName(%s) != %s (actual: %s)", testCase.branchName, testCase.expected, actual)
			}
		})
	}
}
