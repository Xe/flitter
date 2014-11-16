package workflow

// Context represents a context for a workflow process.
type Context struct {
	CleanupTasks []HandleFunc
	Name         string
	Arguments    map[string]string

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
		Name:      name,
		steps:     nil,
		Arguments: make(map[string]string),
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
func (c *Context) Run() (err error) {
	for _, funct := range c.steps {
		err = funct(c)
		if err != nil {
			break
		}
	}

	for _, task := range c.CleanupTasks {
		defer task(c)
	}

	return
}
