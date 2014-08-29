package main

import (
	"github.com/coreos/go-etcd/etcd"
)

// Struct Config holds all the configuration that builder will need in the form of
// strings from etcd. This uses code in another part of the repository to marshall
// and demarshall config in and out of etcd.
type Config struct {
	User   string
	Repo   string
	Branch string
	Sha    string
	etcd   *etcd.Client
}

// NewConfig allocates and retuens a config structure for builder. It also seeds
// the values from etcd.
func NewConfig(uplink string) (c *Config) {
	c = &Config{
		etcd: etcd.NewClient([]string{uplink}),
	}
	return
}
