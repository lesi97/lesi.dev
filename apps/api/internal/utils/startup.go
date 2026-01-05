package utils

import (
	"os"
)

/*
A startup function that should be called on the main thread

that provides a developer with useful information
*/
func Startup(logger *Logger,port string) {
	env := os.Getenv("GO_ENV")
	const protocol = "http://"
	if env == "" {
		env = "development"
	}
	name, _ := readModuleInfo()
	ip := getLocalIp()

	logger.PrintColour(false, "brightWhite", "\n  > %s\n", name) 					// Package name
	logger.PrintColour(false, "brightBlack", "\tEnvironments: .env (%s)", env)		// Current environment
	logger.PrintColour(false, "brightMagenta", "\n\t- Local:")						// Localhost address label
	logger.PrintColour(false, "cyan", "\t  %v%v%v", protocol, "localhost", port)	// Localhost address value
	if ip != "" {
		logger.PrintColour(false, "brightMagenta", "\n\t- Network:")				// Network address label
		logger.PrintColour(false, "cyan", "\t  %v%v%v\n", protocol,ip, port)		// Network address value
	}

	logger.PrintColour(false, "green", "\n  ✓ Server Ready\n\n")					// Server ready
}