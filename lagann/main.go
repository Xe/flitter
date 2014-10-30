package main

import (
	"net/http"
	"os"

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

	routing.HandleFunc("/", root)
	routing.HandleFunc("/register", register)
	routing.HandleFunc("/login", login)

	usermux.Post("/user/create", createApp)
	appmux.Post("/app/candeploy/:app", canDeployApp)
	appmux.Post("/app/deploy/:app", deployApp)

	auth, _ := auth.NewAuth("http://"+os.Getenv("HOST")+":4001", "/flitter/lagann/authkeys/")

	routing.Handle("/user", negroni.New(
		auth,
		negroni.Wrap(usermux),
	))

	routing.Handle("/app", negroni.New(
		negroni.Wrap(appmux),
	))

	n.Run(":3000")
}
