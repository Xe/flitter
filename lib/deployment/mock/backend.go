package mock

import "errors"

// Struct Backend is a mock backend for scheduling and starting Actors.
type Backend struct {
	actors []*Actor
	name   string
}

// Struct BackendInfo is the serializable version of a Backend.
type BackendInfo struct {
	Name       string `json:"name"`
	ActorCount int    `json:"count"`
}

// Info returns the serializable information about a mock Backend.
func (b *Backend) Info() (i *BackendInfo) {
	i = &BackendInfo{
		Name:       "Mock backend",
		ActorCount: len(b.actors),
	}

	return
}

// Deploy introduces a new Actor to the Backend if the Actor is unique
// to the Backend.
func (b *Backend) Deploy(a *Actor) error {
	for _, actor := range b.actors {
		if a.Name() == actor.Name() {
			return errors.New("Duplicate actor")
		}
	}

	b.actors = append(b.actors, a)

	a.Start()

	return nil
}

// ListDeploys returns a list of the Actors in the Backend.
func (b *Backend) ListDeploys(pattern string) (res []*Actor, err error) {
	return b.actors, nil
}

// GetActorer looks for an Actor by name pattern and returns it or
// an error describing the failure.
func (b *Backend) GetActorer(pattern string) (*Actor, error) {
	for _, actor := range b.actors {
		if actor.Name() == pattern {
			return actor, nil
		}
	}

	return nil, errors.New("No such Actor by name " + pattern)
}

// Stop asks the Actor to stop and returns its output.
func (b *Backend) Stop(a *Actor) error {
	return a.Stop()
}

// Start asks the Actor to start and returns its output.
func (b *Backend) Start(a *Actor) error {
	return a.Start()
}

// Restart asks the Actor to restart and returns its output.
func (b *Backend) Restart(a *Actor) error {
	return a.Restart()
}

// Destroy destroys the Actor from the list of active Actors and
// returns an error if the given actor doesn't exist.
func (b *Backend) Destroy(a *Actor) error {
	for i, actor := range b.actors {
		if actor.Name() == a.Name() {
			temp := b.actors[:i]
			temp2 := b.actors[i+1:]
			b.actors = nil

			for _, myactor := range temp {
				b.actors = append(b.actors, myactor)
			}
			for _, myactor := range temp2 {
				b.actors = append(b.actors, myactor)
			}

			return nil
		}
	}

	return errors.New("No such actor loaded")
}
