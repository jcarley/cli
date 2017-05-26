package git

import (
	"io/ioutil"
	"os"
	"os/exec"
	"testing"
)

func TestRm(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current working directory: %s", err)
	}
	defer os.Chdir(wd)

	dir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatalf("Failed to make temp directory: %s", err)
	}

	err = os.Chdir(dir)
	if err != nil {
		t.Fatalf("Failed to change working directory: %s", err)
	}
	err = exec.Command("git", "init").Run()
	if err != nil {
		t.Fatalf("Failed to initialize a git directory: %s", err)
	}

	ig := New()
	remote := "datica"
	err = ig.Add(remote, "git@github.com/github/github.git")
	if err != nil {
		t.Fatalf("Failed to add a git remote: %s", err)
	}

	err = ig.Rm(remote)
	if err != nil {
		t.Fatalf("Failed to remove the git remote: %s", err)
	}

	remotes, err := ig.List()
	if err != nil {
		t.Fatalf("Failed to list git remotes: %s", err)
	}
	if len(remotes) != 0 {
		t.Logf("remotes: %+v %d \"%s\"", remotes, len(remotes), remotes[0])
		t.Fatalf("Unexpected git remotes found: %s", remotes)
	}
}

func TestRmDoesNotExist(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current working directory: %s", err)
	}
	t.Log(wd)
	defer os.Chdir(wd)

	dir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatalf("Failed to make temp directory: %s", err)
	}

	err = os.Chdir(dir)
	if err != nil {
		t.Fatalf("Failed to change working directory: %s", err)
	}
	err = exec.Command("git", "init").Run()
	if err != nil {
		t.Fatalf("Failed to initialize a git directory: %s", err)
	}

	err = New().Rm("datica")
	if err == nil {
		t.Fatalf("Expected an error when removing a non-existant git remote")
	}
}
