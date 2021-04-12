// +build js

package logs

import (
	"os"
)

func IsTerminal(fd int) bool {
	return false
}

func TryToDropFilePageCache(fd int, offset int64, length int64) {
}

func compareFileCreatedTime(a, b os.FileInfo) bool {
	return false
}
