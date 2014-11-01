package main

import (
	"io/ioutil"
	"os/exec"

	"github.com/Xe/flitter/lib/output"
	"github.com/coreos/go-systemd/unit"
)

var unitTemplate = []*unit.UnitOption{}

func startUnit(name, tag string, myunit []*unit.UnitOption) (err error) {
	dir, err := ioutil.TempDir("", "flitter-builder")
	if err != nil {
		return err
	}

	byteslicedunitreader := unit.Serialize(myunit)
	byteslicedunit, err := ioutil.ReadAll(byteslicedunitreader)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(dir+"/"+name+"@.service", byteslicedunit, 0666)
	if err != nil {
		return err
	}

	cmd := exec.Command("fleetctl", "-endpoint", *etcdhost, "start",
		dir+"/"+name+"@"+tag+".service")

	out, err := cmd.CombinedOutput()
	if err != nil {
		output.WriteData(string(out))
		return err
	}

	return err
}
