// Package breakfast contains build-chain logic
package breakfast

import (
	"io/ioutil"

	"github.com/pkg/errors"
	yaml "gopkg.in/yaml.v2"
)

const (
	// FileName is the name of the breakfast file that will be parsed
	FileName = "breakfast.yaml"
)

// File represents the entire breakfast YAML file
type File struct {
	Tasks *Tasks `yaml:"tasks,omitempty"`
}

// Tasks represent the different lifecycle hooks tasks can be executed against
type Tasks struct {
	BeforeBuild []*TaskDeclaration `yaml:"before_build,omitempty"`
}

// TaskDeclaration represents a task declared in the breakfast YAML file
type TaskDeclaration struct {
	Package string                 `yaml:"package,omitempty"`
	Task    string                 `yaml:"task,omitempty"`
	Params  map[string]interface{} `yaml:"params,omitempty"`
}

// Parse parses the yaml file at the given location and returns a File object
func Parse(location string) (*File, error) {
	bytes, err := ioutil.ReadFile(location)
	if err != nil {
		return nil, errors.Wrapf(err, "Error reading breakfast file %s", location)
	}

	file := &File{}
	if err := yaml.Unmarshal(bytes, file); err != nil {
		return nil, errors.Wrapf(err, "Error unmarshaling breakfast file %s", location)
	}

	return file, nil
}
