package main

import (
	"github.com/Xe/flitter/etcdconfig"
	"github.com/coreos/go-etcd/etcd"
)

// Struct Config holds all the configuration that builder will need in the form of
// strings from etcd. This uses code in another part of the repository to marshall
// and demarshall config in and out of etcd.
type Config struct {
	SlugbuilderImage   string `etcd:"/deis/slugbuilder/image"`
	SlugrunnerImage    string `etcd:"/deis/slugrunner/image"`
	DockerfileShim     string `etcd:"/deis/builder/dockerfileshim"`
	ControllerProtocol string `etcd:"/deis/controller/protocol"`
	ControllerHost     string `etcd:"/deis/controller/host"`
	ControllerPort     string `etcd:"/deis/controller/port"`
	ControllerBuildKey string `etcd:"/deis/controller/buildKey"`
	etcd               *etcd.Client
	updates            chan *etcd.Response
	stop               chan bool
}

// NewConfig allocates and retuens a config structure for builder. It also seeds
// the values from etcd.
func NewConfig(prefix, uplink string) (c *Config) {
	c = &Config{
		etcd:    etcd.NewClient([]string{uplink}),
		updates: make(chan *etcd.Response, 10),
		stop:    make(chan bool),
	}

	etcdconfig.Demarshall(c.etcd, c)

	c.etcd.Watch(prefix, 0, true, c.updates, c.stop)

	go func() {
		for update := range c.updates {
			_ = update // TODO: replace me with a better method

			etcdconfig.Demarshall(c.etcd, c)
		}
	}()

	return
}
