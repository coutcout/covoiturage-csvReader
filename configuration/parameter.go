// Used to configure application
package configuration

import (
	"bytes"
	"flag"
	"fmt"
	"os"
)

type Parameters struct {
	ConfigFilePath string
}

func ParseFlag(progname string, args []string) (*Parameters, string, error) {
	flags := flag.NewFlagSet(progname, flag.ContinueOnError)

	var buf bytes.Buffer
	flags.SetOutput(&buf)

	var params Parameters
	flags.StringVar(&params.ConfigFilePath, "config", "./application-dev.yaml", "path to config file")
	
	err := flags.Parse(args)
	if err != nil {
		return nil, buf.String(), err
	}
	
	if err := ValidateConfigPath(params.ConfigFilePath); err != nil {
		return nil, buf.String(), err
	}

	return &params, buf.String(), nil
}

func ValidateConfigPath(path string) error {
	s, err := os.Stat(path)
	if err != nil {
		return err
	}
	if s.IsDir() {
		return fmt.Errorf("'%s' is a directory, not a normal file", path)
	}
	return nil
}