package test

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// SetUpGitRepo runs git init in the current directory.
func SetUpGitRepo() error {
	output, err := RunCommand("git", []string{"init"})
	if err != nil {
		return fmt.Errorf("Unexpected error setting up git repo: %s", output)
	}
	return nil
}

// SetUpAssociation runs the associate command with the appropriate arguments to
// successfully associate to the test environment.
func SetUpAssociation() error {
	output, err := RunCommand(BinaryName, []string{"associate", EnvLabel, SvcLabel, "-a", Alias})
	if err != nil {
		return fmt.Errorf("Unexpected error setting up association: %s", output)
	}
	return nil
}

// ClearAssociations runs the clear --environments command.
func ClearAssociations() error {
	output, err := RunCommand(BinaryName, []string{"clear", "--environments"})
	if err != nil {
		return fmt.Errorf("Unexpected error clearing associations: %s", output)
	}
	return nil
}

// RunCommand runs the given command and arguments with the current os ENV.
func RunCommand(command string, args []string) (string, error) {
	cmd := exec.Command(command, args...)
	cmd.Env = os.Environ()
	cmd.Stdin = strings.NewReader("n\n")
	output, err := cmd.CombinedOutput()
	return string(output), err
}
