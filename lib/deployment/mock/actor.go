package mock

import (
	"fmt"
	"time"

	"github.com/Xe/flitter/lib/deployment"
)

// Struct Actor represents a single worker running a given
// image with given tags.
type Actor struct {
	name   string
	tags   []string
	image  string
	status deployment.Status
}

// Struct ActorInfo is the serializable form of Actor.
type ActorInfo struct {
	Name   string            `json:"name"`
	Tags   []string          `json:"tags"`
	Image  string            `json:"image"`
	Status deployment.Status `json:"status"`
}

// Image returns the Actor's image as a string.
func (a *Actor) Image() string {
	return a.image
}

// Name returns the Actor's name as a string.
func (a *Actor) Name() string {
	return a.name
}

// Tags returns the Actor's metadata tags as a string slice.
func (a *Actor) Tags() []string {
	return a.tags
}

// Status returns the Actor's status as a deployment.Status.
func (a *Actor) Status() deployment.Status {
	return a.status
}

// Info returns the serializable form of the Actor.
func (a *Actor) Info() *ActorInfo {
	return &ActorInfo{
		Name:   a.name,
		Status: a.status,
		Tags:   a.tags,
		Image:  a.image,
	}
}

// Start commands the Actor to start, returning an error
// if the start fails.
func (a *Actor) Start() error {
	a.status = deployment.STATUS_LOADING

	go func() {
		time.Sleep(3 * time.Second)
		a.status = deployment.STATUS_RUNNING
	}()

	return nil
}

// Stop commands the Actor to stop, returning an error if
// the stop fails.
func (a *Actor) Stop() error {
	go func() {
		time.Sleep(3 * time.Second)
		a.status = deployment.STATUS_STOPPED
	}()

	return nil
}

// Restart calls Stop and then Start and returns an error if
// either fails, bailing out on the first failure.
func (a *Actor) Restart() error {
	err := a.Stop()
	if err != nil {
		return err
	}

	// Wait for the actor to stop
	for a.status != deployment.STATUS_STOPPED {
		time.Sleep(1 * time.Second)
	}

	err = a.Start()
	if err != nil {
		return err
	}

	for a.status != deployment.STATUS_RUNNING && a.status == deployment.STATUS_LOADING {
		time.Sleep(1 * time.Second)
	}

	return err
}

// String returns a representation of the Actor as a simple string.
func (a *Actor) String() string {
	return fmt.Sprintf("%s running %s with tags %s and status %d",
		a.name, a.image, a.tags, a.status)
}
