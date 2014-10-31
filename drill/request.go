package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/Xe/flitter/lagann/datatypes"
)

func request(url string, values url.Values) (rep *datatypes.Reply, err error) {
	client := http.Client{}
	rep = &datatypes.Reply{}

	resp, err := client.PostForm(url, values)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(content, rep)
	if err != nil {
		return nil, err
	}

	return
}
