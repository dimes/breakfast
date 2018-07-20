package task

import "context"

// Context in the task input to run a task
type Context struct {
	Ctx         context.Context
	WorkdingDir string
}

// Task represents a build task that can be executed at various stages of a build
type Task interface {
	Execute(ctx *Context) error
}
