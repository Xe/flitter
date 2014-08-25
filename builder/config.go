package main

import (
	"github.com/Xe/flitter/deis"
	"github.com/Xe/flitter/etcdconfig"
	"github.com/coreos/go-etcd/etcd"
)

// Struct Config holds all the configuration that builder will need in the form of
// strings from etcd. This uses code in another part of the repository to marshall
// and demarshall config in and out of etcd.
type Config struct {
	SlugbuilderImage string `etcd:"/deis/slugbuilder/image"`
	SlugrunnerImage  string `etcd:"/deis/slugrunner/image"`
	DockerfileShim   string `etcd:"/deis/builder/dockerfileshim"`
	Controller       *deis.Controller
	etcd             *etcd.Client
	updates          chan *etcd.Response
	stop             chan bool
}

// NewConfig allocates and retuens a config structure for builder. It also seeds
// the values from etcd.
func NewConfig(uplink string) (c *Config) {
	c = &Config{
		etcd:       etcd.NewClient([]string{uplink}),
		updates:    make(chan *etcd.Response, 10),
		stop:       make(chan bool),
		Controller: deis.NewControllerEtcd(uplink),
	}

	etcdconfig.Demarshal(c.etcd, c)

	return
}
