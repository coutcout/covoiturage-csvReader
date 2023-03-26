// Package configuration is used to configure application
package configuration

import (
	"os"

	"gopkg.in/yaml.v3"
)

// Config struct define all available application configurations
type Config struct {
	Server struct {
		Host string `yaml:"host"`
		Port string `yaml:"port"`
	}

	Journey struct {
		Import struct {
			MaxUploadFile int64 `yaml:"max-upload-file-size"`
		}

		Parser struct {
			WorkerPoolSize int `yaml:"worker-pool-size"`
		}
	}
}


// NewConfig creates a new Config from a file. The file must be a YAML config file
// 
// @param configPath - Path to the configuration file
func NewConfig(configPath string)(*Config, error){
	config := &Config{}

	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	d := yaml.NewDecoder(file)
	
	if err := d.Decode(&config); err != nil {
		return nil, err
	}

	return config, nil
}