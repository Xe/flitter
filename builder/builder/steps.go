package main

import (
	"bufio"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/Xe/flitter/lagann/datatypes"
	"github.com/Xe/flitter/lib/output"
	"github.com/Xe/flitter/lib/workflow"
	"github.com/coreos/go-systemd/unit"
)

func makeTempDir(c *workflow.Context) (err error) {
	dir, err := ioutil.TempDir("", "flitter-builder")
	if err != nil {
		return err
	}

	c.Arguments["tempdir"] = dir

	c.CleanupTasks = append(c.CleanupTasks, func(c *workflow.Context) error {
		err = os.RemoveAll(dir)
		if err != nil {
			return err
		}

		return nil
	})

	return
}

func extractTarball(c *workflow.Context) (err error) {
	// Extract branch to deploy
	output.WriteHeader("Extracting " + repo)

	dir, ok := c.Arguments["tempdir"]
	if !ok {
		return errors.New("Impossible state")
	}

	cmd := exec.Command("git", "archive", branch)

	fout, err := os.Create(dir + "/app.tar")
	if err != nil {
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
				return err
			}

			if err != nil {
				output.WriteData(err)
				return err
			}

			output.WriteData(string(line))
		}
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

	return nil
}

func checkDockerfile(c *workflow.Context) (err error) {
	dir, ok := c.Arguments["tempdir"]
	if !ok {
		return errors.New("Impossible state")
	}

	// Find the Dockerfile
	if _, err := os.Stat(dir + "/Dockerfile"); os.IsNotExist(err) {
		output.WriteError("Need a dockerfile to build!")

		output.WriteData("See https://github.com/Xe/flitter/issues/24 for more information")
		return err
	}

	return
}

func injectLayers(c *workflow.Context) (err error) {
	dir, ok := c.Arguments["tempdir"]
	if !ok {
		return errors.New("Impossible state")
	}

	// Inject some vars
	output.WriteHeader("Injecting flitter layers to Dockerfile")

	dockerfout, err := os.OpenFile(dir+"/Dockerfile", os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		output.WriteError("Could not inject things to Dockerfile")
		return err
	}

	_, err = dockerfout.Write([]byte("\nENV GIT_SHA " + sha + "\n"))
	if err != nil {
		output.WriteError("Error: " + err.Error())
		return err
	}

	dockerfout.Write([]byte("ENV BUILTBY " + os.Getenv("USER") + "\n"))
	dockerfout.Write([]byte("ENV APPNAME " + os.Getenv("REPO") + "\n"))
	dockerfout.Write([]byte("ENV BUILDID " + buildid))
	dockerfout.Close()

	output.WriteData("done")

	return
}

func validateDockerfile(c *workflow.Context) (err error) {
	dir, ok := c.Arguments["tempdir"]
	if !ok {
		return errors.New("Impossible state")
	}

	// Validate Docker image
	output.WriteHeader("Validating Dockerfile")

	fin, err := os.Open(dir + "/Dockerfile")
	if err != nil {
		output.WriteError("Could not validate Dockerfile")
		return err
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
				output.WriteData("Please make sure to only expose one port")
				output.WriteData("You can and will run into undefined behavior")
				output.WriteData("You have been warned")
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

	return
}

func buildImage(c *workflow.Context) (err error) {
	dir, ok := c.Arguments["tempdir"]
	if !ok {
		return errors.New("Impossible state")
	}

	// Build docker image
	image = config.RegistryHost + ":" + config.RegistryPort +
		"/" + os.Getenv("USER") + "/" + repo + ":" + sha[:7]
	// 192.168.45.117:5000/xena/mpd:1fc8018

	output.WriteHeader("Building docker image " + image)
	cmd := exec.Command("docker", "build", "-t", image, dir)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout

	err = cmd.Run()
	if err != nil {
		output.WriteError(err.Error())
		return err
	}

	return
}

func tagAndPushImage(c *workflow.Context) (err error) {
	// Tag and push to registry
	output.WriteHeader("Pushing image " + image + " to registry")
	cmd := exec.Command("docker", "push", image)
	err = cmd.Run()
	if err != nil {
		output.WriteError("Error in push: " + err.Error())
		stderr, err := cmd.StderrPipe()
		if err != nil {
			output.WriteData("Cannot get debug information")
			return err
		}

		spew := bufio.NewReader(stderr)

		for {
			line, _, err := spew.ReadLine()

			if err == io.EOF {
				return err
			}

			if err != nil {
				output.WriteData(err)
				return err
			}

			output.WriteData(string(line))
		}
	}

	output.WriteData("done")

	return
}

func deployImage(c *workflow.Context) (err error) {
	// Report information about the build
	output.WriteHeader("Deploying")

	build := &datatypes.Build{
		App:   os.Getenv("REPO"),
		ID:    buildid,
		Image: image,
		User:  os.Getenv("USER"),
	}

	unitSlice := []*unit.UnitOption{
		{"Unit", "Description", "Flitter app " + repo + " deploy " + build.ID},
		{"Service", "TimeoutStartSec", "30m"},
		{"Service", "ExecStartPre", "/usr/bin/docker pull " + build.Image},
		{"Service", "ExecStartPre", "-/usr/bin/docker rm -f app-" + repo + "-" + build.ID},
		{"Service", "ExecStart", "/bin/sh -c '/usr/bin/docker run -P --name app-" + repo + "-" + build.ID + " --hostname " + repo + " -e HOST=$COREOS_PRIVATE_IPV4 " + build.Image + " '"},
		{"Service", "ExecStop", "/usr/bin/docker rm -f app-" + repo + "-" + build.ID},
	}

	if err := startUnit("app-"+repo, buildid, unitSlice); err != nil {
		output.WriteError("Fleet unit start failed: " + err.Error())
		output.WriteData("Please verify that fleet is online.")
		return err
	}

	return
}

func successMessage(c *workflow.Context) (err error) {
	// Print end message
	output.WriteHeader("Success")
	output.WriteData("Your app is in the docker registry as " + image)
	output.WriteData("Your build id is " + buildid)
	output.WriteData("")
	output.WriteData("You may access your app at http://" + repo + "." + config.Domain + " once it spins up")

	return
}
