// Used to configure application
package configuration

import (
	"os"

	"gopkg.in/yaml.v3"
)

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