package main

import (
	"log"
	"os"

	"github.com/codegangsta/cli"
)

var Commands = []cli.Command{
	commandLogin,
	commandLogout,
	commandCreate,
	commandRegister,
	commandWhoami,
}

var commandLogin = cli.Command{
	Name:   "login",
	Usage:  "login to the given controller",
	Action: doLogin,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:   "username",
			Usage:  "specify an arbitrary username instead of getting it over stdin",
			EnvVar: "DRILL_AUTH_USER",
		},
		cli.StringFlag{
			Name:   "password",
			Usage:  "specify an arbitrary password instead of getting it over stdin",
			EnvVar: "DRILL_AUTH_PASSWORD",
		},
	},
}

var commandLogout = cli.Command{
	Name:   "logout",
	Usage:  "log out of the controller",
	Action: doLogout,
}

var commandCreate = cli.Command{
	Name:  "create",
	Usage: "create a new application in the controller",
	Description: `if there is no argument to this command, an application name
will be generated, otherwise the first argument will be the
application name`,
	Action: doCreate,
}

var commandRegister = cli.Command{
	Name:   "register",
	Usage:  "register a new account with the given controller",
	Action: doRegister,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:   "username",
			Usage:  "specify an arbitrary username instead of getting it over stdin",
			EnvVar: "DRILL_AUTH_USER",
		},
		cli.StringFlag{
			Name:   "password",
			Usage:  "specify an arbitrary password instead of getting it over stdin",
			EnvVar: "DRILL_AUTH_PASSWORD",
		},
	},
}

var commandWhoami = cli.Command{
	Name:   "whoami",
	Usage:  "whoami queries the controller to see who it thinks you are",
	Action: doWhoami,
}

func debug(v ...interface{}) {
	if os.Getenv("DEBUG") != "" {
		log.Println(v...)
	}
}

func assert(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func doLogin(c *cli.Context) {
}

func doLogout(c *cli.Context) {
}

func doCreate(c *cli.Context) {
}

func doRegister(c *cli.Context) {
}

func doWhoami(c *cli.Context) {
}
