package main

import (
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/Xe/flitter/lib/output"
	"github.com/coreos/go-systemd/unit"
)

var unitTemplate = []*unit.UnitOption{}

func startUnit(name, tag string, myunit []*unit.UnitOption) (err error) {
	dir, err := ioutil.TempDir("", "flitter-builder")
	if err != nil {
		return err
	}

	defer func() {
		os.RemoveAll(dir)
	}()

	byteslicedunitreader := unit.Serialize(myunit)
	byteslicedunit, err := ioutil.ReadAll(byteslicedunitreader)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(dir+"/"+name+"@.service", byteslicedunit, 0666)
	if err != nil {
		return err
	}

	command := []string{"-endpoint", *etcdhost, "start",
		dir + "/" + name + "@" + tag + ".service"}

	output.WriteData("$ fleetctl " + strings.Join(command, " "))
	output.WriteData("Waiting for fleetctl...")

	cmd := exec.Command("/opt/fleet/fleetctl", command...)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Start()
	if err != nil {
		return err
	}

	err = cmd.Wait()
	if err != nil {
		return err
	}

	output.WriteData("done")

	return err
}
