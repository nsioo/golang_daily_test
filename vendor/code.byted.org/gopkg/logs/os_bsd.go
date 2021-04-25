// +build darwin dragonfly freebsd netbsd openbsd

package logs

import (
	"os"
	"syscall"
	"time"
)

const ioctlReadTermios = syscall.TIOCGETA

func TryToDropFilePageCache(fd int, offset int64, length int64) {
}

func compareFileCreatedTime(a, b os.FileInfo) bool {
	stati := a.Sys().(*syscall.Stat_t)
	statj := b.Sys().(*syscall.Stat_t)
	ctimei := time.Unix(stati.Ctimespec.Sec, stati.Ctimespec.Nsec)
	ctimej := time.Unix(statj.Ctimespec.Sec, statj.Ctimespec.Nsec)
	return ctimei.After(ctimej)
}
