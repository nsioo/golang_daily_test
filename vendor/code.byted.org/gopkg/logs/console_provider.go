package logs

import (
	"fmt"
	"os"
	"runtime"
)

type Brush func(string) string

func NewBrush(color string) Brush {
	pre := "\033["
	reset := "\033[0m"
	return func(text string) string {
		return pre + color + "m" + text + reset
	}
}

var colors = []Brush{
	NewBrush("1;34"), // Trace    blue
	NewBrush("1;34"), // Debug    blue
	NewBrush("1;36"), // Info     cyan
	NewBrush("1;32"), // Notice   green
	NewBrush("1;33"), // Warn     yellow
	NewBrush("1;31"), // Error    red
	NewBrush("1;35"), // Fatal    magenta
}

type ConsoleProvider struct {
	level int
	color bool
}

func NewConsoleProvider() *ConsoleProvider {
	p := &ConsoleProvider{}
	p.color = runtime.GOOS != "windows" && IsTerminal(int(os.Stdout.Fd()))
	return p
}

func (cp *ConsoleProvider) Init() error {
	return nil
}

func (cp *ConsoleProvider) SetLevel(l int) {
	cp.level = l
}

func (cp *ConsoleProvider) WriteMsg(msg string, level int) error {
	if level < cp.level {
		return nil
	}
	if cp.color {
		fmt.Fprint(os.Stdout, colors[level](msg))
	} else {
		fmt.Fprint(os.Stdout, msg)
	}
	return nil
}

func (cp *ConsoleProvider) Flush() error {
	return nil
}

func (cp *ConsoleProvider) Destroy() error {
	return nil
}
