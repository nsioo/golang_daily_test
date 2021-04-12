package logs

import (
	"syscall"
	"unsafe"
	"os"
	"time"
)

var kernel32 = syscall.NewLazyDLL("kernel32.dll")
var procGetConsoleMode = kernel32.NewProc("GetConsoleMode")

func IsTerminal(fd int) bool {
	var st uint32
	r, _, e := syscall.Syscall(procGetConsoleMode.Addr(), 2, uintptr(fd), uintptr(unsafe.Pointer(&st)), 0)
	return r != 0 && e == 0
}

func TryToDropFilePageCache(fd int, offset int64, length int64) {
}
func compareFileCreatedTime(a, b os.FileInfo) bool {
	stati := a.Sys().(*syscall.Win32FileAttributeData)
	statj := b.Sys().(*syscall.Win32FileAttributeData)
	ctimei := time.Unix(0, stati.CreationTime.Nanoseconds())
	ctimej := time.Unix(0, statj.CreationTime.Nanoseconds())
	return ctimei.After(ctimej)
}
