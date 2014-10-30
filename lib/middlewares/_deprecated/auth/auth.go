package auth

import (
	"net/http"

	"github.com/Xe/flitter/lib/utils"
	"github.com/coreos/go-etcd/etcd"
)

// Struct Auth is a negroni middleware for checking X-Lagann-Auth for the correct SSH
// key fingerprint.
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
	}

	if user, allowed := utils.CanConnect(a.e, authkey); allowed {
		r.Header.Set("X-Lagann-User", user)

		next(rw, r)
	} else {
		http.Error(rw, `{"code": 401, "message": "No account found", "data": ["No registration found or user was deleted."]}`, http.StatusUnauthorized)
	}
}
