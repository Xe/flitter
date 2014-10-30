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
		utils.Reply(r, w, "Invalid request: "+err.Error(), 500)
		return
	}
	err = json.Unmarshal(body, app)
	if err != nil {
		utils.Reply(r, w, "Invalid request: "+err.Error(), 500)
		return
	}

	if _, err := client.Get(constants.ETCD_APPS+app.Name, false, false); err == nil {
		utils.Reply(r, w, "App "+app.Name+" already exists", 409)
		return
	} else {
		out, err := json.Marshal(app.Users)
		if err != nil {
			utils.Reply(r, w, "Invalid request: "+err.Error(), 500)
			return
		}

		client.Set(constants.ETCD_APPS+app.Name+"/users", string(out), 0)
		client.Set(constants.ETCD_APPS+app.Name+"/name", app.Name, 0)

		utils.Reply(r, w, "App "+app.Name+" created", 200)
	}
}
