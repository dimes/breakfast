package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/dimes/breakfast"
	"github.com/dimes/breakfast/task"
)

func main() {
	flag.Parse()

	logFlags := log.LstdFlags | log.Lshortfile
	// infoLog := log.New(os.Stderr, "[INFO]", logFlags)
	errLog := log.New(os.Stderr, "[ERROR]", logFlags)

	breakfastFile, err := breakfast.Parse(breakfast.FileName)
	if err != nil {
		errLog.Fatalf("Error parsing breakfast file: %+v", err)
	}

	outputDir, err := ioutil.TempDir(os.TempDir(), "breakfast")
	if err != nil {
		errLog.Fatalf("Error creating temp output directory: %+v", err)
	}
	defer os.RemoveAll(outputDir)

	workingDir, err := os.Getwd()
	if err != nil {
		errLog.Fatalf("Error determining working directory: %+v", err)
	}

	ctx := context.Background()
	env := &task.Env{
		WorkdingDir: workingDir,
	}

	for i, t := range breakfastFile.BeforeBuild {
		output := filepath.Join(outputDir, fmt.Sprintf("%d.so", i))
		task, err := task.NewPackageBuilder(t.Package, t.Task).Build(t.Params, output)
		if err != nil {
			errLog.Fatalf("Error getting task for %s: %+v", t.Package, err)
		}

		task.Execute(ctx, env)
	}
}
