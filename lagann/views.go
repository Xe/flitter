package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/Xe/flitter/lagann/constants"
	"github.com/Xe/flitter/lagann/datatypes"
	"github.com/Xe/flitter/lib/utils"
)

func createApp(w http.ResponseWriter, req *http.Request) {
	app := &datatypes.App{}
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		utils.Reply(req, w, "Invalid request: "+err.Error(), 500)
		return
	}
	err = json.Unmarshal(body, app)
	if err != nil {
		utils.Reply(req, w, "Invalid request: "+err.Error(), 500)
		return
	}

	if _, err := client.Get(constants.ETCD_APPS+app.Name, false, false); err == nil {
		utils.Reply(req, w, "App "+app.Name+" already exists", 409)
		return
	} else {
		out, err := json.Marshal([]string{req.Header.Get("X-Lagann-User")})
		if err != nil {
			utils.Reply(req, w, "Invalid request: "+err.Error(), 500)
			return
		}

		client.Set(constants.ETCD_APPS+app.Name+"/users", string(out), 0)
		client.Set(constants.ETCD_APPS+app.Name+"/name", app.Name, 0)

		utils.Reply(req, w, "App "+app.Name+" created", 200)
	}
}

func addSharing(w http.ResponseWriter, req *http.Request) {
	user := req.Header.Get("X-Lagann-User")
	params := req.URL.Query()
	appname := params.Get(":app")

	err := req.ParseForm()
	if err != nil {
		utils.Reply(req, w, "Internal json error", 500, err)
		return
	}

	res, err := client.Get(constants.ETCD_APPS+appname+"/users", false, false)
	if err != nil {
		utils.Reply(req, w, "No such app "+appname, 404)
		return
	}

	var allowedusers []string
	rawusers := res.Node.Value

	err = json.Unmarshal([]byte(rawusers), &allowedusers)
	if err != nil {
		utils.Reply(req, w, "Internal json error", 500, err)
		return
	}

	for _, username := range allowedusers {
		if strings.ToLower(username) == strings.ToLower(user.Name) {
			goto okay
		}
	}

	utils.Reply(req, w, "Not allowed to modify permissions for "+appname, http.StatusUnauthorized)
	return

okay:

	toadd := req.Form.Get("user")
	allowedusers = append(allowedusers, toadd)

	bs, err := json.Marshal(allowedusers)
	if err != nil {
		utils.Reply(req, w, "Internal json error", 500, err)
		return

	}

	str := string(bs)

	err = client.Set(constants.ETCD_APPS+appname+"/users", str, 0)
	if err != nil {
		utils.Reply(req, w, "Internal etcd error", 500, err)
		return
	}
}
