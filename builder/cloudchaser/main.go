package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"strings"

	"github.com/Xe/flitter/lagann/constants"
	"github.com/Xe/flitter/lib/output"
	"github.com/coreos/go-etcd/etcd"
)

var (
	etcdmachine = flag.String("-etcd-machine", "172.17.42.1", "etcd uplink ip")
)

// main is the entry point for cloudchaser, the build sentry.
func main() {
	if len(os.Args) != 3 {
		log.Fatalln("Usage:\n  cloudchaser <revision> <sha>")
	}

	app := os.Getenv("REPO")
	user := os.Getenv("USER")

	client := etcd.NewClient([]string{"http://" + *etcdmachine + ":4001"})

	output.WriteHeader("Checking permission")
	output.WriteData("user:   " + user)
	output.WriteData("app:    " + app)

	var allowedusers []string

	res, err := client.Get(constants.ETCD_APPS+app+"/users", false, false)
	if err != nil {
		output.WriteError("Permissions check failed: " + err.Error())
		output.WriteData("Do you have permission to deploy this app?")
		os.Exit(1)
	}

	rawusers := res.Node.Value

	err = json.Unmarshal([]byte(rawusers), &allowedusers)
	if err != nil {
		output.WriteError("Internal json decoding reply in allowed app users parsing")
		output.WriteData(err.Error())
		return
	}

	for _, username := range allowedusers {
		if strings.ToLower(username) == strings.ToLower(user) {
			goto allowed
		}
	}

	output.WriteError("User is not authorized to make builds")
	output.WriteData("I think you are " + user)
	output.WriteData("Please check the needed permissions and try again later.")

allowed:

	output.WriteData("")
	output.WriteData("Kicking off build")
	output.WriteData("")

	os.Exit(0)
}
