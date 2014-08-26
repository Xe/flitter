package main

import (
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

		log.Println("Created git repo at " + path)
	}

	return
}
