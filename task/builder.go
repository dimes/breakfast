package task

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"plugin"
	"strings"

	"gopkg.in/yaml.v2"

	"github.com/pkg/errors"
)

// Builder builds tasks into a form consumable by breakfast
type Builder interface {
	Build(params map[string]interface{}, output string) (Task, error)
}

// MainBuilder is a task builder that handles building tasks inside a "main" package
type mainBuilder struct {
	packageName string
	taskName    string

	additionalGoPath string
}

// NewMainBuilder returns a MainBuilder for building the given main package
func NewMainBuilder(packageName, taskName string) Builder {
	return newMainBuilder(packageName, taskName, "")
}

func newMainBuilder(packageName, taskName, additionalGoPath string) *mainBuilder {
	return &mainBuilder{
		packageName:      packageName,
		taskName:         taskName,
		additionalGoPath: additionalGoPath,
	}
}

func (m *mainBuilder) Build(params map[string]interface{}, output string) (Task, error) {
	cmd := exec.Command("go", "build", "-buildmode", "plugin", "-o", output, m.packageName)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	if m.additionalGoPath != "" {
		env := os.Environ()
		for i, e := range env {
			if strings.HasPrefix(e, "GOPATH=") {
				env[i] = fmt.Sprintf("%s%s%s", e, string(os.PathListSeparator), m.additionalGoPath)
				break
			}
		}
		cmd.Env = env
	}

	if err := cmd.Run(); err != nil {
		return nil, errors.Wrap(err, "Could not build task")
	}

	plugin, err := plugin.Open(output)
	if err != nil {
		return nil, errors.Wrap(err, "Could not open plugin")
	}

	sym, err := plugin.Lookup(m.taskName)
	if err != nil {
		return nil, errors.Wrap(err, "Error looking up symbol")
	}

	task, ok := sym.(Task)
	if !ok {
		return nil, errors.New("Symbol was not a task")
	}

	if params != nil {
		marshaledParams, err := yaml.Marshal(params)
		if err != nil {
			return nil, errors.Wrap(err, "Could not marshal params")
		}

		if err := yaml.Unmarshal(marshaledParams, task); err != nil {
			return nil, errors.Wrap(err, "Error binding params to task")
		}
	}

	return task, nil
}

// PackageBuilder is a task builder that builds tasks in non-main packages
type packageBuilder struct {
	packageName string
	taskName    string
}

// NewPackageBuilder returns a builder for the given non-main package
func NewPackageBuilder(packageName, taskName string) Builder {
	return &packageBuilder{
		packageName: packageName,
		taskName:    taskName,
	}
}

func (p *packageBuilder) Build(params map[string]interface{}, output string) (Task, error) {
	taskDir, err := ioutil.TempDir(os.TempDir(), p.taskName)
	if err != nil {
		return nil, errors.Wrap(err, "Error creating temp directory for build")
	}
	defer os.RemoveAll(taskDir)

	srcDir := filepath.Join(taskDir, "src")
	if err := os.MkdirAll(srcDir, os.ModeDir|os.ModePerm); err != nil {
		return nil, errors.Wrap(err, "Could not create src directory")
	}

	contents := "" +
		"package main\n" +
		"import a \"" + p.packageName + "\"\n" +
		"var A a." + p.taskName + "\n" +
		"func main() {}"
	mainFile := filepath.Join(srcDir, "main.go")
	if err := ioutil.WriteFile(mainFile, []byte(contents), os.ModePerm); err != nil {
		return nil, errors.Wrap(err, "Error writing main file")
	}

	return newMainBuilder(mainFile, "A", taskDir).Build(params, output)
}
