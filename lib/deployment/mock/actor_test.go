package mock

import (
	"testing"
	"time"

	"github.com/Xe/flitter/lib/deployment"
)

func makeTestStoppedActor() *Actor {
	return &Actor{
		name:   "Foo",
		tags:   []string{"foo=bar"},
		image:  "flitter/mock",
		status: deployment.STATUS_STOPPED,
	}
}

func makeTestRunningActor() *Actor {
	return &Actor{
		status: deployment.STATUS_RUNNING,
		name:   "Foo",
		tags:   []string{"foo=bar"},
		image:  "flitter/mock",
	}
}

func TestMakeActor(t *testing.T) {
	a := makeTestStoppedActor()

	if a.Name() != "Foo" {
		t.Fatalf("Actor name is %s not Foo", a.Name())
	}

	if a.Tags()[0] != "foo=bar" {
		t.Fatalf("Actor tags are %s not []string{\"foo=bar\"}", a.Tags())
	}

	if a.Image() != "flitter/mock" {
		t.Fatalf("Actor image is %s not flitter/mock", a.Image())
	}
}

func TestStartActor(t *testing.T) {
	a := makeTestStoppedActor()

	if err := a.Start(); err != nil {
		t.Fatalf("Error: %s", err.Error())
	}

	for a.Status() != deployment.STATUS_RUNNING && a.Status() == deployment.STATUS_LOADING {
		time.Sleep(1 * time.Second)
	}

	if a.Status() != deployment.STATUS_RUNNING {
		t.Fatalf("Actor is not running")
	}
}

func TestStopActor(t *testing.T) {
	a := makeTestRunningActor()

	if err := a.Stop(); err != nil {
		t.Fatalf("Actor stop failed. %s", err.Error())
	}

	for a.Status() != deployment.STATUS_STOPPED {
		time.Sleep(1 * time.Second)
	}

	if a.Status() != deployment.STATUS_STOPPED {
		t.Fatal("Actor stopped but is not stopped anymore")
	}
}

func TestRestartActor(t *testing.T) {
	a := makeTestRunningActor()

	if err := a.Restart(); err != nil {
		t.Fatal("Actor restart failed")
	}
}
