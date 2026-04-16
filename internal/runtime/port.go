package runtime

import (
	"fmt"
	"net"
)

// IsPortAvailable checks if a local TCP port is free.
func IsPortAvailable(port int) bool {
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return false
	}
	_ = ln.Close()
	return true
}

// FindAvailablePort starts from the given port and finds the next available one.
func FindAvailablePort(startPort int) int {
	port := startPort
	for port < startPort+100 { // Limit search range
		if IsPortAvailable(port) {
			return port
		}
		port++
	}
	return -1 // Return -1 if none found in range
}
