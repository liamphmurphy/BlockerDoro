package main

import (
	"errors"
	"fmt"
	"os"
)

type Config struct {
	OriginalHostsPath string // this value contains the oldest version of the hosts file
	HostsPath         string
	LocalConfigPath   string
	HostsCopyPath     string
	HostsText         []byte
}

// BackupHosts backs up the currently existing hosts file
func (c *Config) BackupHosts() error {
	hostsCopy, err := os.ReadFile(c.HostsPath)
	if err != nil {
		return fmt.Errorf("error reading %s: %w", c.HostsPath, err)
	}
	c.HostsText = hostsCopy

	err = os.WriteFile(c.HostsCopyPath, hostsCopy, 0644)
	if err != nil {
		return fmt.Errorf("error making copy of hosts file: %w", err)
	}

	// we want to save the original hosts file on the first run, so make that copy.
	originalCopy := fmt.Sprintf("%s/hosts.original", c.LocalConfigPath)
	if _, err := os.Stat(originalCopy); errors.Is(err, os.ErrNotExist) {
		err = os.WriteFile(originalCopy, []byte(hostsCopy), 0644)
		if err != nil {
			return err
		}
		fmt.Printf("saved your original hosts file to %s, PLEASE DO NOT DELETE THIS.\n", originalCopy)
	}
	c.OriginalHostsPath = originalCopy

	return nil
}

// WriteHosts writes to the file in path c.HostsPath with the contents of the passed in newHosts value
func (c *Config) WriteHosts(newHosts string, hostPath string) error {
	err := os.WriteFile(hostPath, []byte(newHosts), 0644)
	if err != nil {
		return err
	}

	return nil
}

// GetHosts takes in a path and returns the string of the file contents
func (c *Config) GetHosts(path string) ([]byte, error) {
	text, err := os.ReadFile(path)
	if err != nil {
		return []byte{}, fmt.Errorf("error reading file: %s", err)
	}
	return text, nil
}

// setup takes an already made Config struct, and populates the needed fields for the app to run
func (c *Config) setup() error {
	c.LocalConfigPath = fmt.Sprintf("%s/blockerdoro", c.LocalConfigPath)
	if _, err := os.Stat(c.LocalConfigPath); errors.Is(err, os.ErrNotExist) {
		os.Mkdir(c.LocalConfigPath, 0755)
	}
	c.HostsCopyPath = fmt.Sprintf("%s/hosts.copy", c.LocalConfigPath)

	c.HostsPath = "/etc/hosts"
	// back up hosts file
	err := c.BackupHosts()
	if err != nil {
		return err
	}

	return nil
}
