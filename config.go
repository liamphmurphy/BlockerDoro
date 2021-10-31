package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	configFile      string // path to the TOML config file
	LocalConfigPath string
	Hosts           Hosts
}

// setup will prepare the config structure for blockerdoro
func (c *Config) setup() error {
	c.LocalConfigPath = fmt.Sprintf("%s/blockerdoro", c.LocalConfigPath)
	if _, err := os.Stat(c.LocalConfigPath); errors.Is(err, os.ErrNotExist) {
		os.Mkdir(c.LocalConfigPath, 0755)
	}

	// prepare config.toml file
	c.configFile = fmt.Sprintf("%s/config.toml", c.LocalConfigPath)
	if _, err := os.Stat(c.LocalConfigPath); errors.Is(err, os.ErrNotExist) {
		var username string
		fmt.Println("Detected that this is the first run of this program. Please enter the username of your current OS user:")
		fmt.Scanln(&username)
		err := writeConfigDefaults(c.configFile)
		if err != nil {
			return fmt.Errorf("error creating the user's default config.toml file: %s", err)
		}
		return fmt.Errorf("%s was not created. Please edit this file manually with your desired config values", c.configFile)
	}

	// setup Hosts struct
	var hosts Hosts
	err := hosts.setup(c.LocalConfigPath)
	if err != nil {
		return fmt.Errorf("error setting up hosts struct: %s", err)
	}

	return nil
}

// writeConfigDefaults writes the default.toml file to the user's config.toml file. This is usually run when no config.toml file for the user exists.
func writeConfigDefaults(path string, username string) error {
	viper.ReadInConfig()

	err = WriteFile(string(defaultValues), path)
	if err != nil {
		return err
	}

	return nil
}
