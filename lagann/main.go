package main

import (
	"flag"
	"net/http"

	"github.com/Xe/flitter/lagann/constants"
	"github.com/Xe/flitter/lib/middlewares/auth"
	"github.com/codegangsta/negroni"
	"github.com/coreos/go-etcd/etcd"
	"github.com/drone/routes"
	"gopkg.in/unrolled/render.v1"
)

var (
	r      *render.Render
	client *etcd.Client

	etcdmachine = flag.String("etcd-machine", "172.17.42.1", "uplink to connect to for etcd")
)

func init() {
	r = render.New(render.Options{})
}

// main is the entry point for lagann
func main() {
	flag.Parse()

	client = etcd.NewClient([]string{"http://" + *etcdmachine + ":4001"})

	routing := http.NewServeMux()
	usermux := routes.New()

	n := negroni.Classic()
	n.UseHandler(routing)

	routing.HandleFunc(constants.ROOT_MUXPATH, root)
	routing.HandleFunc(constants.REGISTER_URL, register)
	routing.HandleFunc(constants.LOGIN_URL, login)

	usermux.Post(constants.APP_CREATE_URL, createApp)
	usermux.Post(constants.APP_SHARING_URL, addSharing)

	auth, _ := auth.NewAuth("http://"+*etcdmachine+":4001", constants.ETCD_LAGANN_AUTHKEYS)

	routing.Handle(constants.USER_MUXPATH, negroni.New(
		auth,
		negroni.Wrap(usermux),
	))

	n.Run(":3000")
}
