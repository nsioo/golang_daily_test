package logs

import (
	"bytes"
	"io"
	"code.byted.org/gopkg/logs/utils"
)

// KVEncoder .
type KVEncoder interface {
	io.Writer
	Reset()
	AppendKVs(kvs ...interface{})
	EndRecord()
	Bytes() []byte
	String() string
}

// TTLogKVEncoder .
// Key(VLen)=Val
//  Name(3)=zyj
type TTLogKVEncoder struct {
	buf *bytes.Buffer
}

// NewTTLogKVEncoder .
func NewTTLogKVEncoder() *TTLogKVEncoder {
	return &TTLogKVEncoder{
		buf: new(bytes.Buffer),
	}
}

func (tte *TTLogKVEncoder) Write(p []byte) (n int, err error) {
	return tte.buf.Write(p)
}

// Reset .
func (tte *TTLogKVEncoder) Reset() {
	tte.buf.Reset()
}

// Bytes .
func (tte *TTLogKVEncoder) Bytes() []byte {
	return tte.buf.Bytes()
}

// String .
func (tte *TTLogKVEncoder) String() string {
	return tte.buf.String()
}

// AppendKVs .
func (tte *TTLogKVEncoder) AppendKVs(kvs ...interface{}) {
	for i := 0; i+1 < len(kvs); i += 2 {
		k := kvs[i]
		v := kvs[i+1]
		kbytes := []byte(utils.Value2Str(k))
		vbytes := []byte(utils.Value2Str(v))
		tte.buf.Write(kbytes)
		tte.buf.Write(equalBytes)
		tte.buf.Write(vbytes)
		tte.buf.Write(spaceBytes)
	}
}

// EndRecord .
func (tte *TTLogKVEncoder) EndRecord() {
	tte.buf.Write(newlineBytes)
}

var (
	lBracketBytes = []byte("(")
	rBracketBytes = []byte(")")
	equalBytes    = []byte("=")
	nilBytes      = []byte("nil")
	newlineBytes  = []byte("\n")
)
