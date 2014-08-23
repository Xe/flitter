package main

import "github.com/coreos/go-etcd/etcd"

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
