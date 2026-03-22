package git

import (
	"os/exec"
	"strings"
)

func GetGitEmail() (string, error) {
	email, err := exec.Command("git", "config", "user.email").Output()
	if err != nil {
		return "", err
	}
	return strings.Trim(string(email), "\n"), nil
}
