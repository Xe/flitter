package main

import (
	"fmt"
	"os"

	"github.com/Xe/flitter/builder/output"
	"github.com/docopt/docopt-go"
)

func main() {
	usage := `Flitter Heroku-ish Slug Builder

Usage:
  builder [options] <repo> <branch>

Options:
  --etcd-prefix=<prefix> Sets the etcd prefix to monitor to <prefix> or
                         [default: /deis]
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

	//_ = config

	if !magicHasPermissionForApp(user, repo) {
		output.WriteError("No permission for " + repo)
		os.Exit(1)
	}

	output.WriteHeader("Building " + repo)

	// Create temporary directory
	// Extract branch to deploy
	output.WriteHeader("Extracting " + repo)
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
