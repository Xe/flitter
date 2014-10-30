package main

import (
	"net/http"
	"os"

	"github.com/Xe/flitter/lagann/constants"
	"github.com/Xe/flitter/lib/middlewares/auth"
	"github.com/codegangsta/negroni"
	"github.com/coreos/go-etcd/etcd"
	"github.com/drone/routes"
	"gopkg.in/unrolled/render.v1"
)

var r *render.Render
var client *etcd.Client

func init() {
	r = render.New(render.Options{})
	client = etcd.NewClient([]string{"http://" + os.Getenv("HOST") + ":4001"})
}

// main is the entry point for lagann
func main() {
	routing := http.NewServeMux()
	usermux := routes.New()
	appmux := routes.New()

	n := negroni.Classic()
	n.UseHandler(routing)

	routing.HandleFunc(constants.ROOT_MUXPATH, root)
	routing.HandleFunc(constants.REGISTER_URL, register)
	routing.HandleFunc(constants.LOGIN_URL, login)

	usermux.Post(constants.APP_CREATE_URL, createApp)
	appmux.Post(constants.DEPLOY_APP_URL, deployApp)

	auth, _ := auth.NewAuth("http://"+os.Getenv("HOST")+":4001", "/flitter/lagann/authkeys/")

	routing.Handle(constants.USER_MUXPATH, negroni.New(
		auth,
		negroni.Wrap(usermux),
	))

	routing.Handle(constants.APP_MUXPATH, negroni.New(
		negroni.Wrap(appmux),
	))

	n.Run(":3000")
}
