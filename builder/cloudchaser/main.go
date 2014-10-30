package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Xe/flitter/lagann/datatypes"
	"github.com/Xe/flitter/lib/output"
)

var (
	etcdmachine = flag.String("-etcd", "http://172.17.42.1:4001", "etcd uplink")
)

// main is the entry point for cloudchaser, the build sentry.
func main() {
	if len(os.Args) != 3 {
		log.Fatalln("Usage:\n  cloudchaser <revision> <sha>")
	}

	config := NewConfig(*etcdmachine)

	app := os.Getenv("REPO")
	user := os.Getenv("USER")

	output.WriteHeader("Checking permission")
	output.WriteData("user:   " + user)
	output.WriteData("app:    " + app)

	posturl := fmt.Sprintf("http://%s:%s/app/candeploy/%s", config.LagannHost,
		config.LagannPort, app)

	userob := &datatypes.User{
		Name: user,
		SSHKeys: []*datatypes.SSHKey{
			{
				Comment:     "cloudchaser!",
				Key:         os.Getenv("KEY"),
				Fingerprint: os.Getenv("FINGERPRINT"),
			},
		},
	}

	jsondata, err := json.Marshal(userob)
	if err != nil {
		output.WriteError("Json encoding error")
		output.WriteData("Please contact support")
		os.Exit(1)
	}

	buf := bytes.NewBuffer(jsondata)

	resp, err := http.Post(posturl, "application/json", buf)
	if err != nil {
		output.WriteError("Error in pemissions check: " + err.Error())
		output.WriteData("Please make sure you have the requisite permissions to")
		output.WriteData("deploy to this app.")
		os.Exit(1)
	}

	if resp.StatusCode != 200 {
		output.WriteError(fmt.Sprintf("Error: %d", resp.StatusCode))
		output.WriteData(resp.Status)

		if resp.StatusCode == 404 {
			output.WriteData("Please recreate " + app + " and try again")
		} else {
			output.WriteData("Please contact support")
		}

		os.Exit(1)
	}

	output.WriteData("")
	output.WriteData("Kicking off build")
	output.WriteData("")

	os.Exit(0)
}
