package env

import (
	"os"
	"strconv"
)

var tceServicePort = ""
var tceDebugPort = ""

const (
	MaxPort = 65535
	MinPort = 0
)

func init() {
	servicePortStr := os.Getenv("TCE_PRIMARY_PORT")
	if isValidPort(servicePortStr) {
		tceServicePort = servicePortStr
	}
	debugPortStr := os.Getenv("TCE_DEBUG_PORT")
	if isValidPort(debugPortStr) {
		tceDebugPort = debugPortStr
	}
}

func TCEServicePort() string {
	return tceServicePort
}

func TCEDebugPort() string {
	return tceDebugPort
}

func isValidPort(portStr string) bool {
	port, err := strconv.Atoi(portStr)
	return err == nil && port >= MinPort && port <= MaxPort
}
