package main

import (
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

func gitClone(c string) error {
	// Clone the given repository to the given directory
	r, err := git.PlainClone("/opt/", false, &git.CloneOptions{
		URL: "git@bitbucket.org:ovoeng/terraform.git",
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
