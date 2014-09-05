// Package workflow implements a simple in memory workflow executor
package workflow

import (
	"reflect"
)

// Event is a string for type checking purposes
type Event string

// Step is a string for type checking purposes
type Step string

// Method is the name of the Step handler we want to dispatch.
// At this time it is expected that the method takes no parameters
// and returns the event which is the result of processing the
// current steps logic in the workflow
type Method string

const (
	// NEW is the first step
	NEW = Step("new")
	// START gets sent to initialize the workflow
	START = Event("start")
)

// Workflow keeps track of current Step of the workflow,
// a map of Step -> StepHandler and
// a routing map Event -> Step
type Workflow struct {
	Document interface{}
	running  bool
	current  Step
	handlers map[Step]Method
	routes   map[Event]Step
}

// NewWorkflow creates a workflow with decent defaults
func NewWorkflow() *Workflow {
	return &Workflow{current: NEW, handlers: make(map[Step]Method), routes: make(map[Event]Step)}
}

// AddHandler adds a handler method for a step
func (wf *Workflow) AddHandler(s Step, h Method) {
	wf.handlers[s] = h
}

// AddRoute links and event to a step
func (wf *Workflow) AddRoute(fromEvent Event, toStep Step) {
	wf.routes[fromEvent] = toStep
}

// Run executes the workflow for doc from the start state
// until it reaches the end of the flow and there are no more handlers
func (wf *Workflow) Run(doc interface{}, startStep Step) {
	wf.running = true
	wf.Document = doc
	wf.current = startStep
	handler := wf.handlers[wf.current]
	for handler != "" {
		d := reflect.ValueOf(wf.Document)
		m := d.MethodByName(string(handler))
		values := m.Call(nil)
		event := values[0].Interface().(Event)
		Step := wf.routes[event]
		if Step != "" {
			wf.current = Step
		}
		handler = wf.handlers[wf.current]
	}
	wf.running = false
}

// IsRunning will tell you if the workflow is busy
func (wf *Workflow) IsRunning() bool {
	return wf.running
}

// Step returns the current step of the workflow
func (wf *Workflow) Step() Step {
	return wf.current
}
