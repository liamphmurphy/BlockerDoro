package main

import (
	"fmt"
	"os"
	"strings"
	"time"
)

func main() {
	var conf Config
	var err error

	// load config data
	err = conf.setup("./config")
	if err != nil {
		e, ok := err.(*FirstRunError)
		if ok {
			fmt.Fprintf(os.Stdout, e.Error())
			os.Exit(0)
		}
		fmt.Fprintf(os.Stderr, "error setting up blockerdoro: %s\n", err)
		os.Exit(1)
	}

	// setup Hosts struct
	var hosts Hosts
	err = hosts.setup("./backups")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error setting up blockerdoro: %s\n", err)
		os.Exit(1)
	}

	conf.Hosts = hosts

	// generate list of domains to block
	newHosts, err := GetNewHostsFile(strings.Join(conf.Domains, "\n"))
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
		time.Sleep(time.Minute * time.Duration(conf.WorkTimer))

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
		time.Sleep(time.Minute * time.Duration(conf.BreakTimer))
	}
}
