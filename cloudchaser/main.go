package main

import (
	"log"
	"os"

	"github.com/Xe/flitter/builder/output"
	"github.com/docopt/docopt-go"
)

func main() {
	usage := `Cloudchaser Builder Sentry

Usage:
  cloudchaser [options] <revision> <sha>

Options:
`
	arguments, err := docopt.Parse(usage, nil, true, "Flitter Sentry 0.1", false)
	if err != nil {
		log.Fatal(err)
	}
	_ = arguments

	output.WriteHeader("Receiver\n")

	output.WriteHeader("Environment:")

	for _, val := range os.Environ() {
		output.WriteData(val)
	}
}
