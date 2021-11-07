package main

import (
	"errors"
	"fmt"
	"os"
)

type Hosts struct {
	OriginalPath string // this value contains the oldest version of the hosts file
	Path         string
	CopyPath     string
	BackupPath   string
	Text         []byte
}

// BackupHosts backs up the currently existing hosts file
func (h *Hosts) BackupHosts(path string) error {
	hostsCopy, err := os.ReadFile(h.Path)
	if err != nil {
		return fmt.Errorf("error reading %s: %w", h.Path, err)
	}
	h.Text = hostsCopy

	err = os.WriteFile(h.CopyPath, hostsCopy, 0644)
	if err != nil {
		return fmt.Errorf("error making copy of hosts file: %w", err)
	}

	return nil
}

// setup takes an already made Config struct, and populates the needed fields for the app to run
func (h *Hosts) setup(backupDir string) error {
	// make the designated backup directory if it does not exist
	if _, err := os.Stat(backupDir); errors.Is(err, os.ErrNotExist) {
		err = os.Mkdir(backupDir, 0644)
		if err != nil {
			return err
		}
	}

	h.CopyPath = fmt.Sprintf("%s/hosts.copy", backupDir)
	h.Path = "/etc/hosts"

	// back up hosts file
	err := h.BackupHosts(h.CopyPath)
	if err != nil {
		return err
	}

	return nil
}

// ReadFile takes in a path and returns the string of the file contents
func ReadFile(path string) ([]byte, error) {
	text, err := os.ReadFile(path)
	if err != nil {
		return []byte{}, fmt.Errorf("error reading file: %s", err)
	}
	return text, nil
}

// WriteFile acts as a general method for writing what is in "contents" to a file.
func WriteFile(contents string, path string) error {
	err := os.WriteFile(path, []byte(contents), 0644)
	if err != nil {
		return err
	}

	return nil
}
