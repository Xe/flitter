package auth

import (
	"net/http"

	"github.com/coreos/go-etcd/etcd"
)

// Struct Auth is a negroni middleware for checking X-Lagann-Auth for the correct uuid
type Auth struct {
	e          *etcd.Client
	BasePath   string
	HeaderName string
}

func NewAuth(machine, path string) (*Auth, error) {
	a := &Auth{
		e:          etcd.NewClient([]string{machine}),
		BasePath:   path,
		HeaderName: "X-Lagann-Auth",
	}

	return a, nil
}

func (a *Auth) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	authkey := r.Header.Get(a.HeaderName)
	if authkey == "" {
		http.Error(rw, `{"code": 401, "message": "Not authorized", "data": ["Try setting `+a.HeaderName+`."]}`, http.StatusUnauthorized)
		return
	}

	node, err := a.e.Get("/flitter/lagann/authkeys/"+authkey, false, false)
	if err != nil {
		http.Error(rw, `{"code": 401, "message": "Not authorized", "data": ["No user found for that API key."]}`, http.StatusUnauthorized)
		return
	}

	user := node.Node.Value
	r.Header.Set("X-Lagann-User", user)

	next(rw, r)
}
