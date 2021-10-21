package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/patrickhuber/shim/pkg/models"
	"github.com/spf13/afero"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v2"
)

func Init(c *cli.Context) error {
	fs := c.App.Metadata["fs"].(afero.Fs)
	source := c.String("source")
	name := c.String("name")
	format := c.String("format")
	return InitInternal(fs, source, name, format)
}

func InitInternal(fs afero.Fs, source, name, format string) error {
	model := &models.Run{
		Name: name,
	}
	var data []byte
	var err error

	switch format {
	case "yaml":
		data, err = yaml.Marshal(model)
	case "json":
		data, err = json.Marshal(model)
	default:
		return fmt.Errorf("unrecognized format %s", format)
	}

	if err != nil {
		return err
	}

	path := strings.TrimSuffix(source, ".exe")
	return afero.WriteFile(fs, path+".yml", data, 0644)
}
