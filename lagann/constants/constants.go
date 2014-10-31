/*
Package constants is a set of constants for Lagann. Right now this is just for
etcd and http routing paths, but in the future it will be more.
*/
package constants

const (
	REGISTER_URL       = "/register"
	LOGIN_URL          = "/login"
	APP_CREATE_URL     = "/user/create"
	CAN_DEPLOY_APP_URL = "/app/candeploy/:app"
	DEPLOY_APP_URL     = "/app/deploy/:app"

	ROOT_MUXPATH = "/"
	USER_MUXPATH = "/user/"
	APP_MUXPATH  = "/app/"

	ETCD_LAGANN_AUTHKEYS = "/flitter/lagann/authkeys/"
	ETCD_LAGANN_USERS    = "/flitter/lagann/users/"
	ETCD_APPS            = "/flitter/apps/"
	ETCD_BUILDER_USERS   = "/flitter/builder/users/"
)
