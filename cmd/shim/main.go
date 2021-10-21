package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/patrickhuber/shim/pkg/models"
	"github.com/patrickhuber/shim/pkg/patch"
	"github.com/spf13/afero"
	"gopkg.in/yaml.v2"
)

func main() {
	handle(run())
}

func show() (string, error) {
	// get currently executing file
	path, err := os.Executable()
	if err != nil {
		return "", err
	}

	fs := afero.NewOsFs()
	file, err := fs.Open(path)
	if err != nil {
		return "", err
	}

	p, err := patch.Get(file)
	if err != nil {
		return "", err
	}
	if p == nil {
		return "", nil
	}
	return string(p.Data), nil
}

func run() error {

	// if someone drops a config file in, it takes priority over the patched shim exe
	runCmd, err := getRunCmdFromConfig()
	if err != nil {
		return err
	}

	// if the config file is missing, try to get from the patch
	if runCmd == nil || len(runCmd.Name) == 0 {
		runCmd, err = getRunCmdFromPatch()
		if err != nil {
			return err
		}
	}

	if runCmd == nil || len(runCmd.Name) == 0 {
		return fmt.Errorf("unable to find shim patch or config file")
	}

	// run the target with the current args
	cmd := exec.Command(runCmd.Name, os.Args[1:]...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin

	return cmd.Run()
}

func getRunCmdFromConfig() (*models.Run, error) {
	path, err := os.Executable()
	if err != nil {
		return nil, err
	}

	path = strings.TrimSuffix(path, ".exe")

	fs := afero.NewOsFs()

	ymlExists, err := afero.Exists(fs, path+".yml")
	if err != nil {
		return nil, err
	}

	jsonExists, err := afero.Exists(fs, path+".json")
	if err != nil {
		return nil, err
	}

	if ymlExists {
		return readYaml(fs, path)
	} else if jsonExists {
		return readJson(fs, path)
	}

	return nil, nil
}

func getRunCmdFromPatch() (*models.Run, error) {
	data, err := show()
	if err != nil {
		return nil, err
	}
	var runCmd *models.Run
	err = yaml.Unmarshal([]byte(data), runCmd)
	if err != nil {
		return nil, err
	}
	return runCmd, nil
}

func readYaml(fs afero.Fs, path string) (*models.Run, error) {

	runCmd := &models.Run{}

	bytes, err := afero.ReadFile(fs, path+".yml")
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(bytes, runCmd)
	if err != nil {
		return nil, err
	}

	return runCmd, nil
}

func readJson(fs afero.Fs, path string) (*models.Run, error) {
	runCmd := &models.Run{}

	bytes, err := afero.ReadFile(fs, path+".json")
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(bytes, runCmd)
	if err != nil {
		return nil, err
	}

	return runCmd, nil
}

func handle(err error) {
	if err == nil {
		return
	}
	switch v := err.(type) {
	case *exec.ExitError:
		errorCode := v.ExitCode()
		os.Exit(errorCode)
	}
	fmt.Fprint(os.Stderr, err.Error())
	os.Exit(1)
}
