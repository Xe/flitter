package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Xe/flitter/builder/output"
)

func main() {
	if len(os.Args) != 3 {
		log.Fatalln("Usage:\n  cloudchaser <revision> <sha>")
	}

	revision := os.Args[1]
	sha := os.Args[2]

	output.WriteHeader("Receiver")

	output.WriteHeader("Environment:")

	for _, val := range os.Environ() {
		output.WriteData(val)
	}

	output.WriteHeader("Kicking off build")
	output.WriteData(fmt.Sprintf("Revision: %s", revision))
	output.WriteData(fmt.Sprintf("Sha: %s", sha))
}
