/*
Package workflow implements a simple workflow.

A workflow consists of a set of States, StateHandlers and Events

The package sets up some defaults :

	NEW   = State("new")
	START = Event("start")


Create a new workflow as follows :

	wf := NewWorkflow()


*/
package workflow
