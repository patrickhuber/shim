package main

import (
	"io"
	"os"

	"github.com/patrickhuber/shim/pkg/patch"
	"github.com/spf13/afero"
	"github.com/urfave/cli/v2"
)

func Apply(c *cli.Context) error {
	fs := c.App.Metadata["fs"].(afero.Fs)
	source := c.String("source")
	target := c.String("target")
	data := c.String("data")

	return ApplyInternal(fs, source, target, data)
}

func ApplyInternal(fs afero.Fs, source, target, data string) error {
	sourceFile, err := fs.Open(source)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	p, err := patch.Get(sourceFile)
	if err != nil {
		return err
	}

	// fetch where to finish writing the original file
	sourceStat, err := sourceFile.Stat()
	if err != nil {
		return err
	}
	sourceSize := sourceStat.Size()
	if p != nil {
		sourceSize -= int64(p.Size)
	}

	// open the data file
	dataBytes, err := afero.ReadFile(fs, data)
	if err != nil {
		return err
	}

	// open the target and truncate for newly writing
	targetFile, err := fs.OpenFile(target, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer targetFile.Close()

	// copy data from source to target
	_, err = io.CopyN(targetFile, sourceFile, sourceSize)
	if err != nil {
		return err
	}

	// create the patch from the datafile
	// apply to the target file
	p = patch.FromString(string(dataBytes))
	return p.Apply(targetFile)
}
