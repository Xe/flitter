package workflow

import "log"

// Context represents a context for a workflow process.
type Context struct {
	CleanupTasks []HandleFunc
	Name         string
	Arguments    map[string]string
	Parameters   map[string]interface{}

	steps []HandleFunc
}

// HandleFunc represents a function to handle a workflow step.
type HandleFunc func(*Context) error

// NilHandleFunc is a stub for a generic handler.
func NilHandleFunc(c *Context) error {
	return nil
}

// New creates and returns a new workflow Context.
func New(name string) (c *Context) {
	c = &Context{
		Name:       name,
		steps:      nil,
		Arguments:  make(map[string]string),
		Parameters: make(map[string]interface{}),
	}

	return
}

// Use adds any number of HandleFuncs to the workflow Context.
func (c *Context) Use(h ...HandleFunc) {
	// prepend h to c.steps
	c.steps = append(h, c.steps...)
}

// Run iteratively runs the HandlerFunc stack until there is no more to do.
// After it finishes it will run the cleanup tasks. If an error happens during
// the job run, it will bail and run all cleanup tasks, then return the error.
// If a panic happens it will perform all cleanup tasks then continue the panic.
func (c *Context) Run() (err error) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Panic! %#v", r)

			c.Cleanup()
			panic(r)
		}
	}()

	defer c.Cleanup()

	for _, funct := range c.steps {
		err = funct(c)
		if err != nil {
			return
		}
	}

	return
}

// Cleanup runs all cleanup tasks in the order they were added. This will eat
// errors and go with its process.
func (c *Context) Cleanup() {
	for _, task := range c.CleanupTasks {
		defer task(c)
	}
}
