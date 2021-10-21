package main

import (
	"os"

	"github.com/spf13/afero"
	"github.com/urfave/cli/v2"
)

func main() {
	fs := afero.NewOsFs()
	path, err := os.Executable()
	if err != nil {
		os.Exit(1)
	}
	// create the patch
	app := &cli.App{
		Metadata: map[string]interface{}{
			"fs":     fs,
			"stdout": os.Stdout,
			"path":   path,
		},
		Commands: []*cli.Command{
			{
				Name:        "apply",
				Aliases:     []string{"a"},
				Description: "Applies the specified patch to the given file",
				Action:      Apply,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "source",
						Aliases:  []string{"s"},
						Required: true,
					},
					&cli.StringFlag{
						Name:     "target",
						Aliases:  []string{"t"},
						Required: true,
					},
					&cli.StringFlag{
						Name:     "data",
						Aliases:  []string{"d"},
						Required: true,
					},
				},
			},
			{
				Name:        "get",
				Aliases:     []string{"g"},
				Description: "Reads the patch from the given file",
				Action:      get,
			},
			{
				Name:        "stat",
				Aliases:     []string{"s"},
				Description: "Reads the patch stats from the given file",
				Action:      Stat,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "source",
						Aliases:  []string{"s"},
						Required: true,
					},
				},
			},
			{
				Name:        "remove",
				Aliases:     []string{"r"},
				Description: "Removes the patch from the given file",
				Action:      remove,
			},
			{
				Name:        "init",
				Aliases:     []string{"i"},
				Description: "Creates a patch file with the given name",
				Action:      Init,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "source",
						Aliases:  []string{"s"},
						Required: true,
					},
					&cli.StringFlag{
						Name:     "name",
						Aliases:  []string{"n"},
						Required: true,
					},
					&cli.StringFlag{
						Name:     "format",
						Aliases:  []string{"f"},
						Required: false,
						Value:    "yaml",
					},
				},
			},
		},
	}
	app.Run(os.Args)
}

func get(c *cli.Context) error {
	return nil
}

func remove(c *cli.Context) error {
	return nil
}
