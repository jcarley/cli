package git

import (
	"os/exec"
	"strings"
)

// List returns a list of all git removes in the current working directory.
func (g *SGit) List() ([]string, error) {
	out, err := exec.Command("git", "remote").Output()
	if err != nil {
		return nil, err
	}
	remotes := []string{}
	for _, r := range strings.Split(string(out), "\n") {
		if len(strings.TrimSpace(r)) > 0 {
			remotes = append(remotes, r)
		}
	}
	return remotes, nil
}
