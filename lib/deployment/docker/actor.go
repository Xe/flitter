package docker

import "github.com/Xe/flitter/lib/deployment"

// Actor is a single unit of work in a docker machine. This represents a
// single container.
type Actor struct {
	name   string
	image  string
	tags   []string
	status deployment.Status
}

// ActorInfo is a parseable version of Actor.
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
