package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"code.google.com/p/go-uuid/uuid"
	"github.com/Xe/dockerclient"
	"github.com/Xe/flitter/lagann/constants"
	"github.com/Xe/flitter/lagann/datatypes"
	"github.com/Xe/flitter/lib/utils"
)

// canDeployApp is mounted at /app/candeploy/:app
//
// It is for Cloudchaser to establish permission to deploy. This should
// be moved to Cloudchaser proper.
func canDeployApp(w http.ResponseWriter, req *http.Request) {
	params := req.URL.Query()
	appname := params.Get(":app")

	user := &datatypes.User{}
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		utils.Reply(r, w, "Invalid request: "+err.Error(), 500)
		return
	}
	err = json.Unmarshal(body, user)
	if err != nil {
		utils.Reply(r, w, "Invalid request: "+err.Error(), 500)
		return
	}

	var allowedusers []string

	// Get app allowed users
	res, err := client.Get(constants.ETCD_APPS+appname+"/users", false, false)
	rawusers := res.Node.Value

	err = json.Unmarshal([]byte(rawusers), &allowedusers)
	if err != nil {
		utils.Reply(r, w, "Internal json decoding reply in allowed app users parsing", 500)
		return
	}

	for _, username := range allowedusers {
		if strings.ToLower(username) == strings.ToLower(user.Name) {
			utils.Reply(r, w, username+" is allowed", 200)
			return
		}
	}

	utils.Reply(r, w, "User is not authorized to make builds", 401)
}

// deployApp is mounted at /app/deploy/:app
//
// This call should also be moved to Cloudchaser or maybe the builder directly.
// It is here mainly to prevent things from breaking too much.
func deployApp(w http.ResponseWriter, req *http.Request) {
	params := req.URL.Query()
	appname := params.Get(":app")

	// Get build object
	build := &datatypes.Build{}

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		utils.Reply(r, w, "Invalid request: "+err.Error(), 400)
		return
	}
	err = json.Unmarshal(body, build)
	if err != nil {
		utils.Reply(r, w, "Invalid request: "+err.Error(), 400)
		return
	}

	if _, err := client.Get(constants.ETCD_APPS+appname, false, false); err != nil {
		utils.Reply(r, w, "No such app "+appname, 404)
		return
	} else {
		// Do fleet deploy here

		uid := uuid.New()[0:8]

		client, err := dockerclient.NewDockerClient("unix:///var/run/docker.sock")
		if err != nil {
			utils.Reply(r, w, "Can't make container", 500, err.Error())
			return
		}

		hc := dockerclient.HostConfig{
			PublishAllPorts: true,
		}

		splitimage := strings.Split(build.Image, "/")
		tag := strings.Split(splitimage[len(splitimage)-1], ":")[1]
		image := strings.Split(build.Image, ":"+tag)[0]

		err = client.PullImage(image, tag)
		if err != nil {
			utils.Reply(r, w, "Can't pull image", 500, err.Error())
		}

		id, err := client.CreateContainer(&dockerclient.ContainerConfig{
			Hostname:   appname,
			Image:      build.Image,
			HostConfig: hc,
		}, "app-"+appname+"-"+build.ID+"-"+uid)
		if err != nil {
			utils.Reply(r, w, "Can't make container", 500, err.Error())
			return
		}
		client.StartContainer(id, &hc)

		utils.Reply(r, w, "App "+appname+" deployed", 200)
	}
}
