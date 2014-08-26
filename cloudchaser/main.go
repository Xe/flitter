package main

import (
	"os"

	"github.com/Xe/flitter/builder/output"
)

func main() {
	output.WriteHeader("Receiver\n")

	output.WriteHeader("Environment:")

	for _, val := range os.Environ() {
		output.WriteData(val)
	}
}
