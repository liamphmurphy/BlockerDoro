package main

import (
	"fmt"
	"os"
	"time"
)

func main() {
	var conf Config
	var err error

	// load config data
	err = conf.setup()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error setting up blockerdoro: %s\n", err)
		os.Exit(1)
	}

	// generate list of domains to block
	newHosts, err := GetNewHostsFile(string(conf.Hosts.Text))
	if err != nil {
		fmt.Fprintf(os.Stderr, "error obtaining values for the new domains: %s\n", err)
		os.Exit(1)
	}

	// begin the timer loop
	for {
		// write the new domains to the hosts file
		err = WriteFile(newHosts, conf.Hosts.Path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error writing new hosts file: %s\n", err)
			os.Exit(1)
		}

		fmt.Println("Beginning work period timer")
		time.Sleep(time.Second * 10)

		// get the content of the original hosts file
		oldHosts, err := ReadFile(conf.Hosts.OriginalPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error reading the original hosts file at %s: %s\n", conf.Hosts.OriginalPath, err)
			os.Exit(1)
		}

		// revert the original hosts copy back to the hosts file
		err = WriteFile(string(oldHosts), conf.Hosts.Path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error reverting the original hosts file, this must be fixed manually: %s\n", err)
			os.Exit(1)
		}

		fmt.Println("Beginning break period timer")
		time.Sleep(time.Second * 10)
	}
}
