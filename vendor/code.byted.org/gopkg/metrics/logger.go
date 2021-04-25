package metrics

import "log"

// LogFunc represents a log func for writing err log
type LogFunc func(format string, v ...interface{})

var logfunc LogFunc = log.Printf

// SetLogFunc updates the global log func.
// it is not goroutine safe, should call it at the beginning of your main function
// default: use log.Printf
func SetLogFunc(fn LogFunc) {
	logfunc = fn
}
