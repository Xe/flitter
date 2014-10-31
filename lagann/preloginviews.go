package main

import (
	"encoding/base64"
	"net/http"

	"code.google.com/p/go-uuid/uuid"

	"github.com/Xe/flitter/lagann/constants"
	"github.com/Xe/flitter/lib/utils"
)

// root is the handler for /. It is a simple 404 page.
func root(w http.ResponseWriter, req *http.Request) {
	utils.Reply(r, w, "No method", http.StatusNotFound)
}

func register(w http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		utils.Reply(r, w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}

	err := req.ParseForm()
	if err != nil {
		utils.Reply(r, w, "Could not parse form "+err.Error(), http.StatusInternalServerError)
		return
	}

	password := utils.HashPassword([]byte(req.Form.Get("password")), []byte("lagann"))
	key := req.Form.Get("sshkey")
	uname := req.Form.Get("username")

	fp := utils.GetFingerprint(key)

	if fp == "error" {
		utils.Reply(r, w, "Invalid SSH key", 401)
		return
	}

	if _, err := client.Get(constants.ETCD_LAGANN_USERS+uname, false, false); err == nil {
		utils.Reply(r, w, "User "+uname+" already exists", 409)
	} else {
		client.Set(constants.ETCD_BUILDER_USERS+uname+"/"+fp, key, 0)

		client.Set(constants.ETCD_LAGANN_USERS+uname+"/password",
			base64.StdEncoding.EncodeToString(password), 0)

		authkey := uuid.New()

		client.Set(constants.ETCD_LAGANN_AUTHKEYS+authkey, uname, 0)

		utils.Reply(r, w, "User "+uname+" created.", 200, map[string]interface{}{
			"authkey": authkey,
			"user":    uname,
		})
	}
}

func login(w http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		utils.Reply(r, w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}

	err := req.ParseForm()
	if err != nil || req.Form.Get("password") == "" || req.Form.Get("username") == "" {
		utils.Reply(r, w, "Could not parse form "+err.Error(), http.StatusInternalServerError)
		return
	}

	password := utils.HashPassword([]byte(req.Form.Get("password")), []byte("lagann"))
	username := req.Form.Get("username")

	comppw := base64.StdEncoding.EncodeToString(password)

	storedpwnode, err := client.Get(constants.ETCD_LAGANN_USERS+username+"/password", false, false)
	if err != nil {
		utils.Reply(r, w, "No such user "+username, http.StatusNotFound)
	}

	if storedpwnode.Node.Value == comppw {
		authkey := uuid.New()

		client.Set(constants.ETCD_LAGANN_AUTHKEYS+authkey, username, 0)

		utils.Reply(r, w, "Logged in as "+username+".", 200, map[string]interface{}{
			"authkey": authkey,
			"user":    username,
		})
	}
}
