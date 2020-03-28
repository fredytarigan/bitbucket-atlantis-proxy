package main

import (
	"io"
	"io/ioutil"
	"log"
	"os"

	"golang.org/x/crypto/ssh"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	gitssh "gopkg.in/src-d/go-git.v4/plumbing/transport/ssh"
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

	var r *git.Repository
	var err error

	sshPath := os.Getenv("HOME") + "/.ssh/repokey"
	sshKey, _ := ioutil.ReadFile(sshPath)
	signer, _ := ssh.ParsePrivateKey(sshKey)
	auth := &gitssh.PublicKeys{
		User:   "git",
		Signer: signer,
	}

	// check if the folder already exists
	if _, err := os.Stat(cloneDir); err != nil {
		if os.IsNotExist(err) {

			// Clone the given repository to the given directory
			log.Printf("git clone %s %s --recursive", repoURL, cloneDir)

			r, err = git.PlainClone(cloneDir, false, &git.CloneOptions{
				URL:               repoURL,
				RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
				Auth:              auth,
			})

			return err
		} else {
			log.Printf("Found existing directory %s", cloneDir)
			r, err = git.PlainOpen(cloneDir)
		}
	}

	// checking out to specific commit provided
	log.Printf("checkout to commit : %s", c)
	w, err := r.Worktree()
	if err != nil {
		return err
	}

	err = w.Checkout(&git.CheckoutOptions{
		Hash: plumbing.NewHash(c),
	})

	// ... retrieving the branch being pointed by HEAD
	ref, err := r.Head()
	if err != nil {
		return err
	}

	// ... retrieving the commit object
	commit, err := r.CommitObject(ref.Hash())
	if err != nil {
		return err
	}

	log.Println(commit)

	return nil
}
