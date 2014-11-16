/*
Command builder is the flitter Docker image builder.
*/
package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/Xe/flitter/lib/output"
	"github.com/Xe/flitter/lib/workflow"
)

var (
	etcdhost = flag.String("etcd-host", "http://172.17.42.1:4001", "etcd url to use")
	help     = flag.Bool("help", false, "shows this message")
	config   *Config
	user     string
	repo     string
	branch   string
	sha      string
	buildid  string
	image    string
)

// main is the entrypoint for the builder
func main() {
	flag.Parse()

	if len(flag.Args()) < 3 {
		*help = true
	}

	if *help {
		fmt.Printf("Usage:\n")
		fmt.Printf("  builder [options] <repo> <branch> <sha>\n\n")
		flag.Usage()
		os.Exit(128)
	}

	config = NewConfig(*etcdhost)
	user = os.Getenv("USER")
	repo = flag.Arg(0)
	branch = flag.Arg(1)
	sha = flag.Arg(2)
	buildid = sha[0:8]

	output.WriteHeader("Building " + repo + " branch " + branch + " as " + user)

	c := workflow.New("builder")

	c.Use(
		makeTempDir,
		extractTarball,
		checkDockerfile,
		injectLayers,
		validateDockerfile,
		buildImage,
		tagAndPushImage,
		deployImage,
		successMessage,
	)

	err := c.Run()
	if err != nil {
		output.WriteError(err.Error())

		log.Fatal(err)
	}
}
