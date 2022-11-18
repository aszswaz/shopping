package config

import (
	"flag"
	"os"
	"path"
	"path/filepath"
)

var configFile string

func GetConfigFile() (string, error) {
	var err error

	if configFile != "" {
		return configFile, nil
	}

	flag.StringVar(&configFile, "c", "", "yml configuration file path. (default shopping.xml in the working directory)")
	flag.Parse()

	if configFile == "" {
		if cwd, err := os.Getwd(); err == nil {
			configFile = path.Join(cwd, "shopping.yml")
		}
	} else {
		if configFile, err = filepath.Abs(configFile); err != nil {
			return configFile, err
		}
	}

	return configFile, nil
}
