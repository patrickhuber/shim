package main

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/patrickhuber/shim/pkg/patch"
	"github.com/spf13/afero"
	"github.com/urfave/cli/v2"
)

func Stat(c *cli.Context) error {
	fs := c.App.Metadata["fs"].(afero.Fs)
	stdout := c.App.Metadata["stdout"].(io.Writer)
	source := c.String("source")
	return StatInternal(fs, source, stdout)
}

func StatInternal(fs afero.Fs, source string, writer io.Writer) error {
	sourceFile, err := fs.Open(source)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	p, err := patch.Get(sourceFile)
	if err != nil {
		return err
	}

	if p == nil {
		return fmt.Errorf("the file is not patched")
	}

	bytes, err := json.Marshal(p)
	if err != nil {
		return err
	}

	_, err = writer.Write(bytes)
	return err
}
