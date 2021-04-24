package env

import (
	"os"
)

var (
	podName string
)

func init() {
	podName = os.Getenv("MY_POD_NAME")
	if podName == "" {
		podName = "-"
	}
}

func InTCE() bool {
	return inTCE
}

func PodName() string {
	return podName
}