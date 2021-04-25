// +build linux

package logs

import (
	"os"
	"syscall"
	"time"
)

const ioctlReadTermios = 0x5401 // syscall.TCGETS
const fadviseDontneed = 4

/* defined in linux-4.14/include/uapi/linux/fadvise.h
 * #define POSIX_FADV_DONTNEED 4
 */

func fadvise(fd int, offset int64, length int64, advice int) (err error) {
	_, _, e := syscall.Syscall6(syscall.SYS_FADVISE64, uintptr(fd), uintptr(offset), uintptr(length), uintptr(advice), 0, 0)
	return e
}

func TryToDropFilePageCache(fd int, offset int64, length int64) {
	fadvise(fd, offset, length, fadviseDontneed)
}

func compareFileCreatedTime(a, b os.FileInfo) bool {
	stati := a.Sys().(*syscall.Stat_t)
	statj := b.Sys().(*syscall.Stat_t)
	ctimei := time.Unix(stati.Ctim.Sec, stati.Ctim.Nsec)
	ctimej := time.Unix(statj.Ctim.Sec, statj.Ctim.Nsec)
	return ctimei.After(ctimej)
}
