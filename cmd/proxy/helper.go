package main

import (
	"io"
	"log"
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
	ok, err := IsDirEmpty(cloneDir)

	if err != nil {
		os.Mkdir(cloneDir, 0755)
	}

	if ok {
		// Clone the given repository to the given directory
		log.Printf("Cloning source repositoring from %s to local %s", repoURL, cloneDir)
		_ = exec.Command("git", "clone", "--branch", "master", "--single-branch", repoURL, cloneDir)
	}

	// run a pull
	log.Printf("Running git pull command on %s", cloneDir)
	cmd := exec.Command("git", "pull")
	cmd.Dir = cloneDir

	// checkout to a given commit
	log.Printf("Running git checkout to commitID %s on %s", c, cloneDir)
	cmd = exec.Command("git", "checkout", c)
	cmd.Dir = cloneDir

	return nil

}
