package main

import (
	"bufio"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/Xe/flitter/lib/output"
	"github.com/Xe/flitter/lib/workflow"
)

func makeTempDir(c *workflow.Context) (err error) {
	dir, err := ioutil.TempDir("", "flitter-builder")
	if err != nil {
		output.WriteError("Could not make temporary directory")
		output.WriteData("Please contact your system administrator")

		return err
	}

	c.Arguments["tempdir"] = dir

	c.CleanupTasks = append(c.CleanupTasks, func(c *workflow.Context) error {
		err = os.RemoveAll(dir)
		if err != nil {
			output.WriteError("\n" + err.Error())
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
