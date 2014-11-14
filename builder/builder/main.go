/*
Command builder is the flitter Docker image builder.
*/
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/Xe/flitter/lagann/datatypes"
	"github.com/Xe/flitter/lib/output"
	"github.com/coreos/go-systemd/unit"
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

	c.Use(makeTempDir)
	c.Use(extractTarball)
	c.Use(checkDockerfile)
	c.Use(injectLayers)
	c.Use(validateDockerfile)
	c.Use(buildImage)
	c.Use(tagAndPushImage)
	c.Use(deployImage)
	c.Use(successMessage)

	err := c.Run()
	if err != nil {
		output.WriteError(err.Error())
	}
}
