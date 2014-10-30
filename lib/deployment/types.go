package deployment

import "fmt"

// Type Status is a flagged 64 bit integer for the state of an Actor.
type Status int64

const (
	STATUS_FAILED  = iota // Actor failed running, needs intervention
	STATUS_RUNNING        // Actor is running as normal
	STATUS_LOADING        // Actor is running a loading task
	STATUS_STOPPED        // Actor has been stopped by an administrator
	STATUS_HUNG           // Actor is running but is not replying to health checks
)

/*
Interface Backender is an interface for a backend for Flitter. It by design is
agnostic and implements a few basic calls that all backends that Flitter supports
has.
*/
type Backender interface {
	fmt.Stringer

	Deploy(Actorer) (err error)                                // Deploys a new Actorer to the Backend.
	ListDeploys(pattern string) (deploys []Actorer, err error) // Lists all Actorers on the Backend.
	GetActorer(pattern string) (Actorer, error)                // Get an Actorer matching a pattern

	Info() interface{} // Backend-specific call for information about the backend.

	Stop(Actorer) error    // Arbitrarily stop a given Actorer from running
	Start(Actorer) error   // Arbitrarily start an Actorer
	Restart(Actorer) error // Arbitrarily restart an Actorer
	Destroy(Actorer) error // Destroy an Actorer
}

/*
Interface Actorer represents a single Actor that Flitter has deployed or otherwise
will know about. This is expected to be implemented by a child struct that embeds
this interface.
*/
type Actorer interface {
	fmt.Stringer

	Image() string     // Get docker image name.
	Name() string      // Get name of Actorer.
	Tags() []string    // Get tags used by Actorer.
	Status() Status    // Get status of Actorer.
	Info() interface{} // Get information about Actorer. Backend specific output.

	Start() error   // Command Actorer to start.
	Stop() error    // Command Actorer to stop.
	Restart() error // Command Actorer to stop and then start again.
}
