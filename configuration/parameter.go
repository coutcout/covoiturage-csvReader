// Package configuration is used to configure application
package configuration

import (
	"bytes"
	"flag"
	"fmt"
	"os"
)

// Parameters struct defines all available arguments of the application
type Parameters struct {
	ConfigFilePath string
}

// ParseFlag parses command line arguments. The program name must be passed as the first argument to this function.
// 
// @param progname - the program name to be used for parsing
// @param args - the command line arguments to
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

// ValidateConfigPath checks if the path is a valid config path.
// 
// @param path - The path to check. Must be a normal file or a directory.
// 
// @return An error if the path is invalid nil otherwise. The error is nil if the path is valid and can be read
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