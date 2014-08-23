package main

import (
	"github.com/Xe/flitter/etcdconfig"
	"github.com/coreos/go-etcd/etcd"
)

// Struct KeySet represents the set of keys that a user could authenticate as.
type KeySet struct {
	Keys map[string]string `etcd:"/deis/builder/users"`
}

// GetKeySet returns a new KeySet for checking authentication.
func GetKeySet() (keyset *KeySet) {
	keyset = &KeySet{
		Keys: make(map[string]string),
	}

	client := etcd.NewClient([]string{"http://127.0.0.1:4001"})

	etcdconfig.Demarshal(client, keyset)
	etcdconfig.Subscribe(client, keyset, "/deis/builder/users")

	return
}
