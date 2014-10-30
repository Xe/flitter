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

// NewAuth returns a new Auth instance and an error or nil based on the machine and
// the base path in etcd to look through. A sane default for Flitter is:
//
//    /flitter/lagann/apikeys/
func NewAuth(machine, path string) (*Auth, error) {
	a := &Auth{
		e:          etcd.NewClient([]string{machine}),
		BasePath:   path,
		HeaderName: "X-Lagann-Auth",
	}

	return a, nil
}

// This middleware will set the X-Lagann-User to the correct value out of etcd
// based on the API key in X-Lagann-Auth.
func (a *Auth) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	authkey := r.Header.Get(a.HeaderName)
	if authkey == "" {
		http.Error(rw, `{"code": 401, "message": "Not authorized", "data": ["Try setting `+a.HeaderName+`."]}`, http.StatusUnauthorized)
		return
	}

	node, err := a.e.Get(a.BasePath+authkey, false, false)
	if err != nil {
		http.Error(rw, `{"code": 401, "message": "Not authorized", "data": ["No user found for that API key."]}`, http.StatusUnauthorized)
		return
	}

	user := node.Node.Value
	r.Header.Set("X-Lagann-User", user)

	next(rw, r)
}
