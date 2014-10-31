package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"strings"

	"github.com/Xe/flitter/lagann/constants"
	"github.com/codegangsta/cli"
	"github.com/howeyc/gopass"
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
		cli.StringFlag{
			Name:   "sshkey",
			Usage:  "ssh key to register with controller",
			Value:  os.Getenv("HOME") + "/.ssh/id_rsa.pub",
			EnvVar: "DRILL_AUTH_KEYPATH",
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
	if len(c.Args()) != 1 {
		fmt.Println("Please specify a controller URL.")
		os.Exit(1)
	}

	controller := c.Args()[0]

	username := c.String("username")
	password := c.String("password")

	if !strings.HasPrefix(controller, "http") {
		fmt.Println(`Please include the full http(|s):// in the URL for the contoller`)
		os.Exit(1)
	}

	var err error

	if username == "" {
		scanner := bufio.NewScanner(os.Stdin)

		fmt.Print("Username: ")
		for {
			scanner.Scan()
			username = scanner.Text()
			if scanner.Err() == nil {
				break
			}

			fmt.Print("Username: ")
		}
	}

	if password == "" {
		fmt.Print("Password: ")
		for {
			password = string(gopass.GetPasswdMasked())
			if password != "" {
				break
			}
			fmt.Print("Password: ")
		}
	}

	values := url.Values{}
	values.Add("username", username)
	values.Add("password", password)

	rep, err := request(controller+constants.LOGIN_URL, values)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		fmt.Printf("Is the controller online?\n")
		os.Exit(1)
	}

	switch rep.Code {
	case 200:
		fmt.Println(rep.Message)

		authkey := rep.Data[0]["authkey"].(string)

		config := &Config{
			URL:     controller,
			AuthKey: authkey,
		}

		err = saveConfig(c.GlobalString("config"), config)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
	case 404:
		fmt.Printf("No such user %s", username)
		os.Exit(1)
	default:
		fmt.Printf("Unknown controller reply %d\n  %s\n", rep.Code, rep.Message)
	}
}

func doRegister(c *cli.Context) {
	if len(c.Args()) != 1 {
		fmt.Println("Please specify a controller URL.")
		os.Exit(1)
	}

	controller := c.Args()[0]

	username := c.String("username")
	password := c.String("password")

	if !strings.HasPrefix(controller, "http") {
		fmt.Println(`Please include the full http(|s):// in the URL for the contoller`)
		os.Exit(1)
	}

	keyfiledata, err := ioutil.ReadFile(c.String("sshkey"))
	if err != nil {
		fmt.Printf("SSH key at %s is unreadable\n", c.String("sshkey"))
		os.Exit(1)
	}

	sshkey := string(keyfiledata)
	sshkey = strings.Split(sshkey, " ")[1]

	if username == "" {
		scanner := bufio.NewScanner(os.Stdin)

		fmt.Print("Username: ")
		for {
			scanner.Scan()
			username = scanner.Text()
			if scanner.Err() == nil {
				break
			}

			fmt.Print("Username: ")
		}
	}

	if password == "" {
		fmt.Print("Password: ")
		for {
			password = string(gopass.GetPasswdMasked())
			if password != "" {
				fmt.Print("Confirm: ")
				otherpass := string(gopass.GetPasswdMasked())
				if password == otherpass {
					break
				} else {
					fmt.Println("Passwords do not match")
				}
			}
			fmt.Print("Password: ")
		}
	}

	values := url.Values{}
	values.Add("username", username)
	values.Add("password", password)
	values.Add("sshkey", sshkey)

	rep, err := request(controller+constants.REGISTER_URL, values)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		fmt.Printf("Is the controller online?\n")
		os.Exit(1)
	}

	fmt.Println(rep.Message)

	switch rep.Code {
	case 200:
		authkey := rep.Data[0]["authkey"].(string)

		config := &Config{
			URL:     controller,
			AuthKey: authkey,
		}

		err = saveConfig(c.GlobalString("config"), config)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
	case 409:
		fmt.Printf("User %s already registered. Try `drill login %s`.", username, controller)
		os.Exit(1)
	default:
		fmt.Printf("Unknown controller reply %d\n  %s\n", rep.Code, rep.Message)
	}
}

func doLogout(c *cli.Context) {
	os.Remove(c.GlobalString("config"))

	fmt.Println("Logged out.")
}

func doCreate(c *cli.Context) {
	fmt.Println("Not implemented")
	os.Exit(1)
}

func doWhoami(c *cli.Context) {
	fmt.Println("Not implemented")
	os.Exit(1)
}
