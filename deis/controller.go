package deis

import (
	"fmt"

	"github.com/Xe/flitter/etcdconfig"
	"github.com/coreos/go-etcd/etcd"
)

// Struct Controller represents the base data structure for the Deis controller.
type Controller struct {
	Protocol string `etcd:"/deis/controller/protocol"`
	Host     string `etcd:"/deis/controller/host"`
	Port     string `etcd:"/deis/controller/port"`
	BuildKey string `etcd:"/deis/controller/builderKey"`
}

// NewController allocates and returns a new controller object with an arbitrary
// hostname and port.
func NewController(host, port string) (c *Controller) {
	c = &Controller{
		Host: host,
		Port: port,
	}

	return
}

// NewControllerEtcd allocates and returns a new controller object with the host
// and port seeded from etcd.
func NewControllerEtcd(etcdUplink string) (c *Controller) {
	c = &Controller{}

	etcdconfig.Demarshal(etcd.NewClient([]string{etcdUplink}), c)

	return
}

// GetURL returns the URL path to the Deis controller.
func (c *Controller) GetURL() (url string) {
	url = fmt.Sprintf("%s://%s:%s", c.Protocol, c.Host, c.Port)

	return
}

// String satisfies the fmt.Stringer interface. It calls GetURL.
func (c *Controller) String() string {
	return c.GetURL()
}
