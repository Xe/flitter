package main

import (
	"encoding/base64"
	"net/http"

	"code.google.com/p/go-uuid/uuid"

	"github.com/Xe/flitter/lagann/constants"
	"github.com/Xe/flitter/lagann/datatypes"
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

	user := &datatypes.User{
		Name:     req.Form.Get("username"),
		Password: string(password),
		SSHKeys: []*datatypes.SSHKey{
			&datatypes.SSHKey{
				Key:         req.Form.Get("sshkey"),
				Fingerprint: utils.GetFingerprint(req.Form.Get("sshkey")),
			},
		},
	}

	if _, err := client.Get(constants.ETCD_LAGANN_USERS+user.Name, false, false); err == nil {
		utils.Reply(r, w, "User "+user.Name+" already exists", 409)
	} else {
		for _, key := range user.SSHKeys {
			client.Set(constants.ETCD_BUILDER_USERS+user.Name+"/"+key.Fingerprint, key.Key, 0)
		}

		client.Set(constants.ETCD_LAGANN_USERS+user.Name+"/password",
			base64.StdEncoding.EncodeToString(password), 0)

		authkey := uuid.New()

		client.Set(constants.ETCD_LAGANN_AUTHKEYS+authkey, user.Name, 0)

		utils.Reply(r, w, "User "+user.Name+" created.", 200, map[string]interface{}{
			"authkey": authkey,
			"user":    user.Name,
		})
	}
}

func login(w http.ResponseWriter, req *http.Request) {
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
