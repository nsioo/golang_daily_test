package utils

import (
	"sync"
	"unicode/utf8"
)

// Use simple []byte instead of bytes.Buffer to avoid large dependency.
type buffer []byte

func (b *buffer) Write(p []byte) {
	*b = append(*b, p...)
}

func (b *buffer) WriteString(s string) {
	*b = append(*b, s...)
}

func (b *buffer) WriteByte(c byte) error {
	*b = append(*b, c)
	return nil
}

func (bp *buffer) WriteRune(r rune) {
	if r < utf8.RuneSelf {
		*bp = append(*bp, byte(r))
		return
	}

	b := *bp
	n := len(b)
	for n+utf8.UTFMax > cap(b) {
		b = append(b, 0)
	}
	w := utf8.EncodeRune(b[n:n+utf8.UTFMax], r)
	*bp = b[:n+w]
}

// pp is used to store a printer's state and is reused with sync.Pool to avoid allocations.
type pp struct {
	buf buffer
}

var ppFree = sync.Pool{
	New: func() interface{} {
		return new(pp)
	},
}

func newPrinter() *pp {
	p := ppFree.Get().(*pp)
	return p
}

func (p *pp) free() {
	p.buf = p.buf[:0]
	ppFree.Put(p)
}

func (p *pp) doPrintf(format string, a []string) {
	end := len(format)
	argNum := 0
	for i := 0; i < end; {
		lasti := i
		for i < end && format[i] != '%' {
			i++
		}
		if i > lasti {
			p.buf.WriteString(format[lasti:i])
		}
		if i >= end {
			// done processing format string
			break
		}

		// Process one verb
		i++
		if i >= end {
			// done processing format string
			break
		}

		c := format[i]
		switch c {
		case 's':
			if argNum < len(a) {
				p.printArgString(a[argNum])
			}
			argNum++
		default:
			p.printByte(c)
		}

		i++
	}

	if argNum < len(a) {
		for _, arg := range a[argNum:] {
			p.printArgString(arg)
		}
	}
}

func (p *pp) printByte(c byte) {
	p.buf.WriteByte(c)
}

func (p *pp) printArgString(a string) {
	p.buf.WriteString(a)
}

func Sprintf(format string, a ...string) string {
	p := newPrinter()
	p.doPrintf(format, a)
	s := string(p.buf)
	p.free()
	return s
}
