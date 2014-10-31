package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/coreos/go-systemd/unit"
)

var unitTemplate = []*unit.UnitOption{}

func startUnit(name string, myunit []*unit.UnitOption) (err error) {
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
	}

	return
}
