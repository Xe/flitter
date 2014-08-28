package main

import (
	"github.com/coreos/go-etcd/etcd"
)

// Struct Config holds all the configuration that builder will need in the form of
// strings from etcd. This uses code in another part of the repository to marshall
// and demarshall config in and out of etcd.
type Config struct {
	SlugbuilderImage string `etcd:"/deis/slugbuilder/image"`
	SlugrunnerImage  string `etcd:"/deis/slugrunner/image"`
	DockerfileShim   string `etcd:"/deis/builder/dockerfileshim"`
	User             string
	Repo             string
	Branch           string
	Sha              string
	etcd             *etcd.Client
	updates          chan *etcd.Response
	stop             chan bool
}

// NewConfig allocates and retuens a config structure for builder. It also seeds
// the values from etcd.
func NewConfig(uplink string) (c *Config) {
	c = &Config{
		etcd: etcd.NewClient([]string{uplink}),
	}

	node, err := c.etcd.Get("/deis/slugbuilder/image", false, false)
	if err != nil {
		c.SlugbuilderImage = "deis/slugbuilder"
	} else {
		c.SlugrunnerImage = node.Node.Value
	}

	node, err = c.etcd.Get("/deis/slugrunner/image", false, false)
	if err != nil {
		c.SlugrunnerImage = "deis/slugrunner"
	} else {
		c.SlugrunnerImage = node.Node.Value
	}

	node, err = c.etcd.Get("/deis/builder/dockerfileshim", false, false)
	if err != nil {
		c.DockerfileShim = `FROM ` + c.SlugbuilderImage + `
RUN mkdir -p /app
WORKDIR /app
ENTRYPOINT ["/runner/init"]
ADD slug.tgz /app
`
	} else {
		c.SlugrunnerImage = node.Node.Value
	}

	return
}
