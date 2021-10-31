package main

import (
	"fmt"
	"strings"
)

var (
	domains    = []string{"www.reddit.com", "www.youtube.com"}
	redirectIP = "0.0.0.0"
)

// GetNewHostsFile takes the original text from the copied hosts file, and edits it to include the domains we want to block.
// That new list is returned as a []string, this function will not perform the file overwrite
func GetNewHostsFile(newDomains string) (string, error) {
	if len(newDomains) == 0 {
		return "", fmt.Errorf("blockerdoro seems to think the original hosts file was empty")
	}

	hosts := strings.Split(string(newDomains), "\n")
	for _, domain := range domains {
		hosts = append(hosts, fmt.Sprintf("%s\t%s", redirectIP, domain))
	}

	return strings.Join(hosts, "\n"), nil
}
