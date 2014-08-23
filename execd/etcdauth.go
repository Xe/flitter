package main

import (
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"io"
	"strings"

	"github.com/coreos/go-etcd/etcd"
)

func CanConnect(e *etcd.Client, user, key string) bool {
	reply, err := e.Get("/deis/builder/users/"+user, false, false)
	if err != nil {
		return false
	}

	for _, node := range reply.Node.Nodes {
		if node.Value == key {
			return true
		}
	}

	return false
}

func getFingerprint(key string) string {
	key = strings.Split(key, " ")[1]
	h := md5.New()

	data, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		return key
	}

	io.WriteString(h, string(data))
	return addColons(fmt.Sprintf("%x", h.Sum(nil)))
}

func addColons(s string) (r string) {
	if len(s) == 0 {
		return ""
	}

	for i, c := range s {
		r = r + string(c)
		if i%2 == 1 && i != len(s)-1 { // Even number, add colon
			r = r + ":"
		}
	}

	return
}
