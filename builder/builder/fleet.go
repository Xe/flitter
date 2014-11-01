package main

import (
	"io/ioutil"
	"os/exec"

	"github.com/Xe/flitter/lib/output"
	"github.com/coreos/go-systemd/unit"
)

var unitTemplate = []*unit.UnitOption{}

func startUnit(name string, myunit []*unit.UnitOption) (err error) {
	dir, err := ioutil.TempDir("", "flitter-builder")
	if err != nil {
		return err
	}

	byteslicedunitreader := unit.Serialize(myunit)
	byteslicedunit, err := ioutil.ReadAll(byteslicedunitreader)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(dir+"/"+name, byteslicedunit, 0666)
	if err != nil {
		return err
	}

	cmd := exec.Command("fleetctl", "--endpoint", *etcdhost, "start",
		dir+"/"+name)

	out, err := cmd.CombinedOutput()
	if err != nil {
		output.WriteData(string(out))
		return err
	}

	return err
}
