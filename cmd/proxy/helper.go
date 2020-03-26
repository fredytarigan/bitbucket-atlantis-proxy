package main

import (
	"io"
	"os"
	"os/exec"
)

func IsDirEmpty(name string) (bool, error) {
	f, err := os.Open(name)
	if err != nil {
		return false, err
	}
	defer f.Close()

	// read in ONLY one file
	_, err = f.Readdir(1)

	// and if the file is EOF... well, the dir is empty.
	if err == io.EOF {
		return true, nil
	}
	return false, err
}

func gitClone(c string) error {
	repoURL := "git@bitbucket.org:ovoeng/terraform.git"
	cloneDir := "/opt/terraform"

	// check if directory exists and not empty
	ok, err := IsDirEmpty("/opt/terraform")

	if err != nil {
		return err
	}

	if ok {
		// Clone the given repository to the given directory
		_ = exec.Command("git", "clone", "--branch", "master", "--single-branch", repoURL, cloneDir)
	}

	// run a pull
	cmd := exec.Command("git", "pull")
	cmd.Dir = cloneDir

	// checkout to a given commit
	cmd = exec.Command("git", "checkout", c)
	cmd.Dir = cloneDir

	return nil

}
