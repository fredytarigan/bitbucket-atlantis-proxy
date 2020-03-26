package main

import (
	"os/user"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/ssh"
)

func gitClone(c string) error {
	currentUser, err := user.Current()

	if err != nil {
		return err
	}

	sshAuth, err := ssh.NewPublicKeysFromFile("git", currentUser.HomeDir+"/.ssh/id_rsa", "")

	if err != nil {
		return err
	}

	// Clone the given repository to the given directory
	r, err := git.PlainClone("/opt/terraform", false, &git.CloneOptions{
		URL:               "git@bitbucket.org:ovoeng/terraform.git",
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
		Auth:              sshAuth,
	})

	if err != nil {
		return err
	}

	w, err := r.Worktree()
	// handle error
	if err != nil {
		return err
	}

	// ... checking out to commit
	err = w.Checkout(&git.CheckoutOptions{
		Hash: plumbing.NewHash(c),
	})

	// handle error
	if err != nil {
		return err
	}

	return nil
}
