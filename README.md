# Breakfast

Breakfast is a super simple toolkit for Go builds

## Example

Breakfast executes tasks in order to complete a build. The interface for task is

```
// Context in the task input to run a task
type Context struct {
	Ctx         context.Context
	WorkdingDir string
}

// Task represents a build task that can be executed at various stages of a build
type Task interface {
	Execute(ctx *Context) error
}
```

You can define your own tasks like so

```
type GreetingTask struct {
	Greeting string `yaml:"greeting"`
}

func (g *GreetingTask) Execute(ctx *task.Context) error {
	fmt.Println("Hello from " + g.Greeting)
}
```

You configure Breakfast with a YAML file

before_build:
  - package: github.com/my/greeting/task
    task: GreetingTask
    params:
      greeting: Breakfast

Execute builds using the CLI

    > breakfast
    => Hello from Breakfast
