package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Xe/flitter/builder/output"
	"github.com/docopt/docopt-go"
)

func main() {
	usage := `Cloudchaser Builder Sentry

Usage:
  cloudchaser <command> <revision>

Commands:
  pre   runs precommit checks, exits with 1 if a failure happens
  post  kicks off builder
`
	arguments, err := docopt.Parse(usage, nil, true, "Flitter Sentry 0.1", false)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(arguments)

	if arguments["<command>"].(string) == "pre" {
		output.WriteHeader("Receiver\n")

		output.WriteHeader("Environment:")

		for _, val := range os.Environ() {
			output.WriteData(val)
		}
	}
}
