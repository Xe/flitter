package main

import (
	"log"
	"os"
	"os/exec"
	"sync"
)

var cacheLock sync.Mutex

// makeGitRepo makes a git repository on the disk at a given path. It returns an error if
// the git command fails, or nil if it doesn't.
// Inspiration from https://github.com/flynn/flynn/blob/master/gitreceived/gitreceived.go#L305
func makeGitRepo(path string) (err error) {
	cacheLock.Lock()
	defer cacheLock.Unlock()

	if _, err = os.Stat(path); os.IsNotExist(err) {
		os.MkdirAll(path, 0755)
		cmd := exec.Command("git", "init", "--bare")
		cmd.Dir = path
		err = cmd.Run()

		os.Mkdir(path+"/cache", 777)

		log.Println("Created git repo at " + path)
	}

	return
}
