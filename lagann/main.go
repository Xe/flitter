package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"code.google.com/p/go-uuid/uuid"
	"github.com/Xe/dockerclient"
	"github.com/Xe/flitter/lagann/datatypes"
	"github.com/codegangsta/negroni"
	"github.com/coreos/go-etcd/etcd"
	"github.com/coreos/go-systemd/unit"
	"github.com/drone/routes"
	"gopkg.in/unrolled/render.v1"
)

func main() {
	r := render.New(render.Options{})
	mux := routes.New()

	client := etcd.NewClient([]string{"http://" + os.Getenv("HOST") + ":4001"})

	mux.Get("/", func(w http.ResponseWriter, req *http.Request) {
		reply(r, w, "No method", 500)
	})

	mux.Post("/register", func(w http.ResponseWriter, req *http.Request) {
		user := &datatypes.User{}
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			reply(r, w, "Invalid request: "+err.Error(), 500)
		}
		err = json.Unmarshal(body, user)
		if err != nil {
			reply(r, w, "Invalid request: "+err.Error(), 500)
		}

		if _, err := client.Get("/flitter/builder/users/"+user.Name, false, false); err != nil {
			reply(r, w, "User "+user.Name+" already exists", 409)
		} else {
			for _, key := range user.SSHKeys {
				client.Set("/flitter/builder/users/"+user.Name+"/"+key.Fingerprint, key.Key, 0)
			}

			reply(r, w, "User "+user.Name+" created.", 200)
		}
	})

	mux.Post("/create", func(w http.ResponseWriter, req *http.Request) {
		app := &datatypes.App{}
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			reply(r, w, "Invalid request: "+err.Error(), 500)
			return
		}
		err = json.Unmarshal(body, app)
		if err != nil {
			reply(r, w, "Invalid request: "+err.Error(), 500)
			return
		}

		if _, err := client.Get("/flitter/apps/"+app.Name, false, false); err == nil {
			reply(r, w, "App "+app.Name+" already exists", 409)
			return
		} else {
			out, err := json.Marshal(app.Users)
			if err != nil {
				reply(r, w, "Invalid request: "+err.Error(), 500)
				return
			}
			client.Set("/flitter/apps/"+app.Name+"/users", string(out), 0)
			client.Set("/flitter/apps/"+app.Name+"/name", app.Name, 0)

			reply(r, w, "App "+app.Name+" created", 200)
		}
	})

	mux.Post("/candeploy/:app", func(w http.ResponseWriter, req *http.Request) {
		params := req.URL.Query()
		appname := params.Get(":app")

		user := &datatypes.User{}
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			reply(r, w, "Invalid request: "+err.Error(), 500)
			return
		}
		err = json.Unmarshal(body, user)
		if err != nil {
			reply(r, w, "Invalid request: "+err.Error(), 500)
			return
		}

		var allowedusers []string

		// Get app allowed users
		res, err := client.Get("/flitter/apps/"+appname+"/users", false, false)
		rawusers := res.Node.Value

		err = json.Unmarshal([]byte(rawusers), &allowedusers)
		if err != nil {
			reply(r, w, "Internal json decoding reply in allowed app users parsing", 500)
			return
		}

		for _, username := range allowedusers {
			if strings.ToLower(username) == strings.ToLower(user.Name) {
				reply(r, w, username+" is allowed", 200)
				return
			}
		}

		reply(r, w, "User is not authorized to make builds", 401)
	})

	mux.Post("/deploy/:app", func(w http.ResponseWriter, req *http.Request) {
		params := req.URL.Query()
		appname := params.Get(":app")

		// Get build object
		build := &datatypes.Build{}

		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			reply(r, w, "Invalid request: "+err.Error(), 400)
			return
		}
		err = json.Unmarshal(body, build)
		if err != nil {
			reply(r, w, "Invalid request: "+err.Error(), 400)
			return
		}

		if _, err := client.Get("/flitter/apps/"+appname, false, false); err != nil {
			reply(r, w, "No such app "+appname, 404)
			return
		} else {
			// Do fleet deploy here

			uid := uuid.New()[0:8]

			unitSlice := []*unit.UnitOption{
				{"Unit", "Description", "Flitter app " + appname + " deploy " + uid},
				{"Service", "EnvironmentFile", "/etc/environment"},
				{"Service", "ExecStartPre", "/usr/bin/docker pull " + build.Image},
				{"Service", "ExecStartPre", "-/usr/bin/docker rm -f app-" + appname + "-" + build.ID + "-%i"},
				{"Service", "ExecStart", "/bin/sh -c '/usr/bin/docker run -P --name app-" + appname + "-" + build.ID + "-%i --hostname " + appname + " -e HOST=$COREOS_PRIVATE_IPV4 " + build.Image + " '"},
				{"Service", "ExecStop", "/usr/bin/docker rm -f app-" + appname + "-" + build.ID + "-%i"},
			}

			/*for startUnit("app-"+appname+"@"+uid, myunit) != nil {
				log.Println("Trying to launch app-" + appname + "@" + uid + "...")
			}*/

			client, err := dockerclient.NewDockerClient("unix:///var/run/docker.sock")
			if err != nil {
				reply(r, w, "Can't make container", 500, err.Error())
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
				reply(r, w, "Can't pull image", 500, err.Error())
			}

			id, err := client.CreateContainer(&dockerclient.ContainerConfig{
				Hostname:   appname,
				Image:      build.Image,
				HostConfig: hc,
			}, "app-"+appname+"-"+build.ID+"-"+uid)
			if err != nil {
				reply(r, w, "Can't make container", 500, err.Error())
				return
			}
			client.StartContainer(id, &hc)

			reply(r, w, "App "+appname+" deployed", 200, unitSlice)
		}
	})

	n := negroni.Classic()
	n.UseHandler(mux)
	n.Run(":3000")
}

func reply(r *render.Render, w http.ResponseWriter, message string, code int, data ...interface{}) {
	r.JSON(w, code, map[string]interface{}{
		"message": message,
		"code":    code,
		"data":    data,
	})
}

func startUnit(name string, myunit []*unit.UnitOption) (err error) {
	/*
		outmap := map[string]interface{}{
			"desiredState": "launched",
			"options":      myunit,
		}

		jsonstr, _ := json.Marshal(outmap)
		client := &http.Client{}
		req, err := http.NewRequest("PUT", "unix:///fleet.sock/v1-alpha/units/"+name+".service",
			bytes.NewBuffer(jsonstr))

		resp, err := client.Do(req)

		if resp.StatusCode != 200 {
			err = errors.New(resp.Status)
		}*/
	return
}
