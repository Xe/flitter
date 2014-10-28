package main

import (
	"github.com/Xe/flitter/etcdconfig"
	"github.com/coreos/go-etcd/etcd"
)

// Struct Config holds all the configuration that builder will need in the form of
// strings from etcd. This uses code in another part of the repository to marshall
// and demarshall config in and out of etcd.
type Config struct {
	etcd *etcd.Client

	Registry string `etcd:"/flitter/registry"`
}

// NewConfig allocates and retuens a config structure for builder. It also seeds
// the values from etcd.
func NewConfig(uplink string) (c *Config) {
	c = &Config{
		etcd: etcd.NewClient([]string{uplink}),
	}
	etcdconfig.Demarshal(c.etcd, c)
	return
}
