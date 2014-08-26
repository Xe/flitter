package main

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"sync"
)

var cacheLock sync.Mutex

// Inspiration from https://github.com/flynn/flynn/blob/master/gitreceived/gitreceived.go#L305
func makeGitRepo(path string) (err error) {
	cacheLock.Lock()
	defer cacheLock.Unlock()

	if _, err = os.Stat(path); os.IsNotExist(err) {
		os.MkdirAll(path, 0755)
		cmd := exec.Command("git", "init", "--bare")
		cmd.Dir = path
		err = cmd.Run()
		return err
	}

	log.Println("Created git repo at " + path)

	err = ioutil.WriteFile(path+"/hooks/pre-receive", []byte(`#!/bin/bash

set -eo pipefail; while read oldrev newrev refname; do
	/app/cloudchaser pre $refname
done`), 0755)
	if err == nil {
		return err
	}

	err = ioutil.WriteFile(path+"/hooks/post-receive", []byte(`#!/bin/bash

set -eo pipefail; while read oldrev newrev refname; do
	/app/cloudchaser post $refname
done`), 0755)

	return
}
