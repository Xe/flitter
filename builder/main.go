/*
Command builder is the flitter Heroku-ish slug builder.
*/
package main

import (
	"bufio"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"

	"github.com/Xe/flitter/builder/output"
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

	config.User = user
	config.Repo = repo
	config.Branch = branch
	config.Sha = sha

	curdir, err := os.Getwd()

	output.WriteData("Running in " + curdir)

	output.WriteHeader("Building " + repo + " branch " + branch + " as " + user)

	// Create temporary directory
	dir, err := ioutil.TempDir("", "flitter-builder")
	if err != nil {
		output.WriteError("Could not create temporary directory")
		output.WriteData("Please contact your system administrator")

		os.Exit(1)
	}

	// Extract branch to deploy
	output.WriteHeader("Extracting " + repo + " to " + dir)

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
	}

	fout.Sync()
	fout.Close()

	// Extract tarball
	cmd = exec.Command("tar", "xf", "app.tar")
	cmd.Dir = dir

	err = cmd.Run()
	if err != nil {
		output.WriteHeader("Error in extracting tarball: " + err.Error())
		os.Exit(1)
	}

	// Grab config from controller
	// Find the Dockerfile or Procfile
	var dockerbuild bool
	if _, err := os.Stat(dir + "/Dockerfile"); os.IsNotExist(err) {
		dockerbuild = false
	} else {
		dockerbuild = true
	}

	if !dockerbuild {
		// Process through slugbuilder if needed
		output.WriteHeader("Building Heroku procfile-based app\n")

		// TODO: add environment from deis controller
		ctidCmd := exec.Command("docker", "run", "-i", "-d", "-v", dir+":/tmp/app", "-v",
			"/home/git"+repo+"/cache"+":/tmp/cache:rw", "deis/slugbuilder")

		ctidBs, err := ctidCmd.Output()
		if err != nil {
			output.WriteError("Error in Heroku build: " + err.Error())
			os.Exit(1)
		}

		ctid := strings.TrimSuffix(string(ctidBs), "\n")

		buildCmd := exec.Command("docker", "attach", ctid)
		buildCmd.Stdout = os.Stdout
		buildCmd.Stderr = os.Stdout

		err = buildCmd.Run()
		if err != nil {
			output.WriteError("Error in Heroku build (attach phase): " + err.Error())
			os.Exit(1)
		}

		fout, err := os.Create(dir + "/Dockerfile")
		if err != nil {
			output.WriteError("Error in Heroku Dockerfile Create: " + err.Error())
			os.Exit(1)
		}

		_, err = fout.Write([]byte(config.DockerfileShim))
		if err != nil {
			output.WriteError("Error in Heroku Dockerfile write: " + err.Error())
		}

		fout.Close()
	}

	// Build docker image
	output.WriteHeader("Building docker image")
	cmd = exec.Command("docker", "build", "-t", repo+":"+sha[:7], dir)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout

	err = cmd.Run()
	if err != nil {
		output.WriteError(err.Error())
		os.Exit(1)
	}

	// Tag and push to registry
	// Extract process types from procfile
	// Report information about the build
	// Print end message
	// Do cleanup of repo and builder
	output.WriteHeader("Cleanup")
	output.WriteData("Removing temporary files")

	err = os.RemoveAll(dir)
	if err != nil {
		output.WriteError("\n" + err.Error())
	}
}
