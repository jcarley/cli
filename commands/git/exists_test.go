package git

import (
	"io/ioutil"
	"os"
	"os/exec"
	"testing"
)

func TestExists(t *testing.T) {
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

	if !New().Exists() {
		t.Fatalf("Exists returned false for a valid git repo")
	}
}

func TestExistsFails(t *testing.T) {
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

	if New().Exists() {
		t.Fatalf("Exists returned true for an empty directory")
	}
}
