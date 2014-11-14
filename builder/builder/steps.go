package main

import (
	"io/ioutil"
	"os"

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
