package main

import (
	"os"

	"github.com/codegangsta/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "drill"
	app.Version = Version
	app.Usage = "a command line frontend to flitter"
	app.Commands = Commands

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "c,config",
			Value:  os.Getenv("HOME") + "/.drill.yaml",
			EnvVar: "DRILL_CONFIG",
			Usage:  "configuration file for drill",
		},
	}

	app.Run(os.Args)
}
