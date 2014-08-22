package main

import (
	"fmt"

	"github.com/docopt/docopt-go"
)

func main() {
	usage := `Flitter Heroku-ish Slug Builder

Usage:
  builder [options] <user> <repo> <branch>

Options:
  --etcd-prefix=<prefix> Sets the etcd prefix to monitor to <prefix> or
                         [default: /deis]
  --etcd-host=<host>     Sets the etcd url to use [default: http://172.17.42.1:4001]
  -h,--help              Show this screen
  -v,--verbose           Show all raw commands as they are running and
                         all output of all commands, even ones that are
                         normally silenced.
  --version              Show version
  --repository-tag=<tag> Tags built docker images with <tag> if set and
                         does not tag them if not.
`

	arguments, _ := docopt.Parse(usage, nil, true, "Flitter Builder 0.1", false)
	fmt.Println(arguments)
}
