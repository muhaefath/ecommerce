package internal

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v2"
)

type Config struct {
	// Environment config
	Env string `yaml:"environment"`
	// Database configuration
	Database map[string]*DatabaseConfig `yaml:"database"`
}

type DatabaseConfig struct {
	Driver      string        `yaml:"driver"`
	Master      string        `yaml:"master"`
	Slaves      []string      `yaml:"slaves"`
	MaxIdleTime time.Duration `yaml:"max_idle_time"`
	MaxLifeTime time.Duration `yaml:"max_life_time"`
	MaxIdleConn int           `yaml:"max_idle_conns"`
	MaxOpenConn int           `yaml:"max_open_conns"`
}

// InitConfig Read and process config file
func InitConfig() Config {
	appconfig := Config{}
	configFilePath := os.Getenv("CONFIG_FILE_PATH")
	if configFilePath == "" {
		// relative path from /build. Use in local
		configFilePath = ".config.yaml"
	}
	if err := readConfig(&appconfig, configFilePath); err != nil {
		fmt.Println("error: ", err)
	}
	return appconfig
}

// readConfig is file handler for reading configuration files into variable
// Return: - boolean
func readConfig(ac *Config, fname string) error {
	fname, err := filepath.Abs(fname)
	if err != nil {
		return err
	}

	yamlFile, err := ioutil.ReadFile(fname) //#nosec
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(yamlFile, ac)
	if err != nil {
		return err
	}

	return nil
}
