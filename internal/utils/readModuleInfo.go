package utils

import (
	"os"
	"strings"
)

/*
Function to read the version number and module name from the `go.mod` file
*/
func readModuleInfo() (name, version string) {
	content, err := os.ReadFile("go.mod")
	if err != nil {
		return "unknown-module", "v0.0.0"
	}

	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "module ") {
			name = strings.TrimSpace(strings.TrimPrefix(line, "module"))
		}
		if strings.HasPrefix(line, "// Version:") {
			version = strings.TrimSpace(strings.TrimPrefix(line, "// Version:"))
		}
	}

	if version == "" {
		version = "v0.0.0"
	}
	return name, version
}