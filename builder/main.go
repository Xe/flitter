/*
Flitter Heroku-ish Slug Builder

Usage:

      builder [options] <repo> <branch>

Options:

      --etcd-host=<host>     Sets the etcd url to use [default: http://172.17.42.1:4001]
      -h,--help              Show this screen
      -v,--verbose           Show all raw commands as they are running and
                             all output of all commands, even ones that are
                             normally silenced.
      --version              Show version
      --repository-tag=<tag> Tags built docker images with <tag> if set and
                             does not tag them if not.

This program assumes it is being run from the root path of all the repositories
that Flitter is tracking.
*/
package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/Xe/flitter/builder/output"
	"github.com/docopt/docopt-go"
)

func main() {
	usage := `Flitter Heroku-ish Slug Builder

Usage:
  builder [options] <repo> <branch>

Options:
  --etcd-host=<host>     Sets the etcd url to use [default: http://172.17.42.1:4001]
  -h,--help              Show this screen
  -v,--verbose           Show all raw commands as they are running and
                         all output of all commands, even ones that are
                         normally silenced.
  --version              Show version
  --repository-tag=<tag> Tags built docker images with <tag> if set and
                         does not tag them if not.

This program assumes it is being run from the root path of all the repositories
that Flitter is tracking.
`

	arguments, _ := docopt.Parse(usage, nil, true, "Flitter Builder 0.1", false)
	fmt.Println(arguments)

	//config := NewConfig(arguments["--etcd-prefix"].(string), arguments["--etcd-host"].(string))
	user := os.Getenv("USER")
	repo := arguments["<repo>"].(string)
	branch := arguments["<branch>"].(string)

	//_ = config

	output.WriteHeader("Building " + repo + " branch " + branch + " as " + user)

	// Create temporary directory
	dir, err := ioutil.TempDir("", "flitter-builder")
	if err != nil {
		output.WriteError("Could not create temporary directory")
		output.WriteData("Please contact your system administrator")

		os.Exit(1)
	}

	// Extract branch to deploy
	output.WriteHeader("Extracting " + repo + " to " + dir)

	cmd := exec.Command("git", "archive", branch)
	cmd.Dir = repo

	fout, err := os.Create(dir + "/app.tar")
	if err != nil {
		output.WriteError("Cannot create application tarball")
		os.Exit(1)
	}

	cmd.Stdout = fout
	cmd.Stderr = os.Stderr

	cmd.Run()

	err = cmd.Wait()
	if err != nil {
		output.WriteError("Git archive problem: " + err.Error())

		stderr, err := cmd.StderrPipe()
		if err != nil {
			output.WriteData("Cannot get debug information")
			os.Exit(1)
		}

		spew := bufio.NewReader(stderr)

		for {
			line, _, err := spew.ReadLine()

			if err == io.EOF {
				os.Exit(1)
			}

			if err != nil {
				output.WriteData(err)
				os.Exit(1)
			}

			output.WriteData(string(line))
		}
	}

	fout.Sync()
	fout.Close()

	// Grab config from controller / etcd
	// Find the Dockerfile or Procfile
	// Process through slugbuilder if needed
	// Build docker image
	// Tag and push to registry
	// Extract process types from procfile
	// Report information about the build
	// Print end message
	// Do cleanup of repo and builder
}
