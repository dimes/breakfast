package task

import "context"

// Env represents the task execution environment
type Env struct {
	WorkdingDir string
}

// Task represents a build task that can be executed at various stages of a build
type Task interface {
	Execute(ctx context.Context, env *Env) error
}
