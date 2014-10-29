/*
Command builder is the flitter Docker image builder.
*/
package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"code.google.com/p/go-uuid/uuid"

	"github.com/Xe/flitter/lagann/datatypes"
	"github.com/Xe/flitter/lib/output"
	"github.com/docopt/docopt-go"
)

func main() {
	usage := `Flitter Image Builder

Usage:
  builder [options] <repo> <branch> <sha>

Options:
  --etcd-host=<host>     Sets the etcd url to use [default: http://172.17.42.1:4001]
  -h,--help              Show this screen
  -v,--verbose           Show all raw commands as they are running and
                         all output of all commands, even ones that are
                         normally silenced.
  --version              Show version
  --repository-tag=<tag> Tags built docker images with <tag> if set and
                         does not tag them if not.

This program assumes it is being run in the bare repository it is building.
`

	arguments, err := docopt.Parse(usage, nil, true, "Flitter Builder 0.1", false)
	if err != nil {
		log.Fatal(err)
	}

	config := NewConfig(arguments["--etcd-host"].(string))
	user := os.Getenv("USER")
	repo := arguments["<repo>"].(string)
	branch := arguments["<branch>"].(string)
	sha := arguments["<sha>"].(string)
	buildid := uuid.New()[0:8]

	output.WriteHeader("Building " + repo + " branch " + branch + " as " + user)

	// Create temporary directory
	dir, err := ioutil.TempDir("", "flitter-builder")
	if err != nil {
		output.WriteError("Could not create temporary directory")
		output.WriteData("Please contact your system administrator")

		os.Exit(1)
	}

	defer func() {
		output.WriteHeader("Cleanup")
		output.WriteData("Removing temporary files")

		err = os.RemoveAll(dir)
		if err != nil {
			output.WriteError("\n" + err.Error())
		}

	}()

	// Extract branch to deploy
	output.WriteHeader("Extracting " + repo)

	cmd := exec.Command("git", "archive", branch)

	fout, err := os.Create(dir + "/app.tar")
	if err != nil {
		output.WriteError("Cannot create application tarball")
		os.Exit(1)
	}

	cmd.Stderr = os.Stderr

	out, err := cmd.Output()
	if err != nil {
		output.WriteHeader("Error in capturing tarball: " + err.Error())

		stderr, err := cmd.StderrPipe()
		if err != nil {
			output.WriteData("Cannot get debug information")
			os.Exit(1)
		}

		spew := bufio.NewReader(stderr)

		for {
			line, _, err := spew.ReadLine()

			if err == io.EOF {
				os.Exit(1)
			}

			if err != nil {
				output.WriteData(err)
				os.Exit(1)
			}

			output.WriteData(string(line))
		}

		os.Exit(1)
	}

	_, err = fout.Write(out)
	if err != nil {
		output.WriteHeader("Error in writing tarball: " + err.Error())
		os.Exit(1)
	}

	fout.Sync()
	fout.Close()

	output.WriteData("done")

	// Extract tarball
	cmd = exec.Command("tar", "xf", "app.tar")
	cmd.Dir = dir

	err = cmd.Run()
	if err != nil {
		output.WriteHeader("Error in extracting tarball: " + err.Error())
		os.Exit(1)
	}

	os.Remove(dir + "/app.tar")

	// Grab config from controller
	// Find the Dockerfile or Procfile
	var dockerbuild bool
	if _, err := os.Stat(dir + "/Dockerfile"); os.IsNotExist(err) {
		dockerbuild = false
		output.WriteError("Need a dockerfile to build!")
		os.Exit(1)
	} else {
		dockerbuild = true
	}

	if !dockerbuild {
		// Process through slugbuilder if needed
		output.WriteHeader("Building Heroku procfile-based app\n")

		// TODO: add environment from deis controller
		ctidCmd := exec.Command("docker", "run", "-idv", dir+":/tmp/app", "-v",
			"/home/git/"+repo+"/cache"+":/tmp/cache:rw", "deis/slugbuilder:latest")

		ctidBs, err := ctidCmd.CombinedOutput()
		if err != nil {
			output.WriteError("Error in Heroku build: " + err.Error())
			output.WriteData(fmt.Sprintf("%s", ctidBs))
			os.Exit(1)
		}

		ctid := strings.TrimSuffix(string(ctidBs), "\n")

		buildCmd := exec.Command("docker", "attach", ctid)
		buildCmd.Stdout = os.Stdout
		buildCmd.Stderr = os.Stderr
		buildCmd.Stdin = strings.NewReader("exit\n")

		err = buildCmd.Run()
		if err != nil {
			output.WriteError("Error in Heroku build (attach phase): " + err.Error())
			os.Exit(1)
		}

		exec.Command("docker", "rm", "-f", ctid)

		fout, err := os.Create(dir + "/Dockerfile")
		if err != nil {
			output.WriteError("Error in Heroku Dockerfile Create: " + err.Error())
			os.Exit(1)
		}

		_, err = fout.Write([]byte(`FROM deis/slugrunner:latest
RUN mkdir -p /app
WORKDIR /app
ENTRYPOINT ["/runner/init"]
ADD slug.tgz /app`))
		if err != nil {
			output.WriteError("Error in Heroku Dockerfile write: " + err.Error())
		}

		fout.Close()
	}

	// Inject some vars
	output.WriteHeader("Injecting flitter layers to Dockerfile")

	dockerfout, err := os.OpenFile(dir+"/Dockerfile", os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal("Could not inject things to Dockerfile")
	}

	_, err = dockerfout.Write([]byte("ENV GIT_SHA " + sha + "\n"))
	if err != nil {
		output.WriteError("Error: " + err.Error())
		os.Exit(1)
	}

	dockerfout.Write([]byte("ENV BUILTBY " + os.Getenv("USER") + "\n"))
	dockerfout.Write([]byte("ENV APPNAME " + os.Getenv("REPO") + "\n"))
	dockerfout.Write([]byte("ENV BUILDID " + buildid))
	dockerfout.Close()

	output.WriteData("done")

	// Validate Docker image
	output.WriteHeader("Validating Dockerfile")

	fin, err := os.Open(dir + "/Dockerfile")
	if err != nil {
		output.WriteError("Could not validate Dockerfile")
		os.Exit(1)
	}

	exposed := false
	scanner := bufio.NewReader(fin)
	line, isPrefix, err := scanner.ReadLine()
	for err == nil && !isPrefix {
		s := string(line)

		split := strings.Split(s, " ")
		if len(split) == 0 {
			continue
		}

		if strings.ToUpper(split[0]) == "EXPOSE" {
			if exposed {
				output.WriteData("Multiple ports exposed")
				output.WriteData("Please make sure to only expose one port.")
				output.WriteData("You can and will run into undefined behavior.")
				output.WriteData("")
				break
			} else {
				exposed = true
			}
		}

		line, isPrefix, err = scanner.ReadLine()
	}

	fin.Close()

	output.WriteData("done")

	// Build docker image
	image := config.RegistryHost + ":" + config.RegistryPort +
		"/" + os.Getenv("USER") + "/" + repo + ":" + sha[:7]
	// 192.168.45.117:5000/xena/mpd:1fc8018

	output.WriteHeader("Building docker image " + image)
	cmd = exec.Command("docker", "build", "-t", image, dir)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout

	err = cmd.Run()
	if err != nil {
		output.WriteError(err.Error())
		os.Exit(1)
	}

	defer func() {
		cmd := exec.Command("docker", "rmi", image)
		cmd.Run()
	}()

	// Tag and push to registry
	output.WriteHeader("Pushing image " + image + " to registry")
	cmd = exec.Command("docker", "push", image)
	err = cmd.Run()
	if err != nil {
		output.WriteError("Error in push: " + err.Error())
		stderr, err := cmd.StderrPipe()
		if err != nil {
			output.WriteData("Cannot get debug information")
			os.Exit(1)
		}

		spew := bufio.NewReader(stderr)

		for {
			line, _, err := spew.ReadLine()

			if err == io.EOF {
				os.Exit(1)
			}

			if err != nil {
				output.WriteData(err)
				os.Exit(1)
			}

			output.WriteData(string(line))
		}

		os.Exit(1)
	}

	output.WriteData("done")

	// Extract process types from procfile
	// Report information about the build
	output.WriteHeader("Deploying")

	build := &datatypes.Build{
		App:   os.Getenv("REPO"),
		ID:    buildid,
		Image: image,
		User:  os.Getenv("USER"),
	}

	jsonstr, _ := json.Marshal(build)

	output.WriteData("Sending build summary to lagann")
	resp, err := http.Post(
		"http://"+config.LagannHost+":"+config.LagannPort+"/deploy/"+os.Getenv("REPO"), "application/json",
		bytes.NewBuffer(jsonstr))
	if err != nil {
		output.WriteError("Error: " + err.Error())
		output.WriteData("Is lagann online?")
		output.WriteData("Skipping deploy")
	} else {
		if resp.StatusCode != 200 {
			output.WriteError("Error: " + resp.Status)
			output.WriteData(fmt.Sprintf("Status code %d", resp.StatusCode))
			os.Exit(1)
		}
	}
	output.WriteData("done")

	// Print end message
	output.WriteHeader("Success")
	output.WriteData("Your app is in the docker registry as " + image)
	output.WriteData("Your build id is " + buildid)
	output.WriteData("Your application will spawn on the next available node")
}
