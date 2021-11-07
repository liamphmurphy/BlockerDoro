package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	configFile string   // path to the TOML config file
	Domains    []string `mapstructure:"domains"`
	WorkTimer  int      `mapstructure:"worktimer"`
	BreakTimer int      `mapstructure:"breaktimer"`
	Hosts      Hosts
}

// FirstRunError usage indicates a special error relating to running BlockerDoro for the first time, not considered fatal.
// For example, this "error" is used when the config.toml file is written for the first time and requires manual editing before BlockerDoro is useful.
type FirstRunError struct {
	Err error
}

func (e *FirstRunError) Error() string {
	return e.Err.Error()
}

// setup will prepare the config structure for blockerdoro. The bool indicates success or failure
func (c *Config) setup(configDir string) error {
	if _, err := os.Stat(configDir); errors.Is(err, os.ErrNotExist) {
		os.Mkdir(configDir, 0755)
	}

	// prepare config.toml file
	c.configFile = fmt.Sprintf("%s/config.toml", configDir)
	if _, err := os.Stat(c.configFile); errors.Is(err, os.ErrNotExist) {
		err := writeConfigDefaults(c.configFile)
		if err != nil {
			return fmt.Errorf("error creating the user's default config.toml file: %s", err)
		}

		return &FirstRunError{
			Err: fmt.Errorf("first time setup detected, please update %s manually\n", c.configFile),
		}
	}

	err := c.populateConfig()
	if err != nil {
		return fmt.Errorf("error populating the config struct: %s", err)
	}

	return nil
}

// popualteConfig will attempt to read in the local config file, and will update the associated Config struct accordingly if successful.
// This probably won't be unit tested, as that would mean just basically testing the viper library, which is out of scope.
func (c *Config) populateConfig() error {
	v := viper.New()
	v.SetConfigFile(c.configFile)
	if err := v.ReadInConfig(); err != nil {
		return err
	}

	err := v.Unmarshal(c)
	if err != nil {
		return err
	}

	return nil
}

// writeConfigDefaults writes the default.toml file to the user's config.toml file. This is usually run when no config.toml file for the user exists.
func writeConfigDefaults(path string) error {
	viper.SetDefault("workTimer", 20)
	viper.SetDefault("breakTimer", 5)
	viper.SetDefault("domains", []string{})

	err := viper.WriteConfigAs(path)
	if err != nil {
		return err
	}

	return nil
}
