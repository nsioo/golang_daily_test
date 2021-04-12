package metrics

import "errors"

const (
	msgpackStringHeaderSize = 3
	msgpackArrayHeaderSize  = 3
)

// https://github.com/msgpack/msgpack/blob/master/spec.md#formats-array
/*
array 16 stores an array whose length is upto (2^16)-1 elements:
+--------+--------+--------+~~~~~~~~~~~~~~~~~+
|  0xdc  |YYYYYYYY|YYYYYYYY|    N objects    |
+--------+--------+--------+~~~~~~~~~~~~~~~~~+
*/
func msgpackAppendArrayHeader(b []byte, n uint16) []byte {
	return append(b, 0xdc, byte(n>>8), byte(n))
}

// https://github.com/msgpack/msgpack/blob/master/spec.md#formats-str
/*
str 16 stores a byte array whose length is upto (2^16)-1 bytes:
+--------+--------+--------+========+
|  0xda  |ZZZZZZZZ|ZZZZZZZZ|  data  |
+--------+--------+--------+========+
*/
func msgpackAppendStringHeader(b []byte, n uint16) []byte {
	return append(b, 0xda, byte(n>>8), byte(n))
}

func msgpackAppendString(b []byte, s string) []byte {
	b = msgpackAppendStringHeader(b, uint16(len(s)))
	return append(b, s...)
}

func msgpackStringSize(s string) int {
	return msgpackStringHeaderSize + len(s)
}

// for test only
func msgpackUnpackArrayHeader(b []byte) (int, []byte, error) {
	if len(b) < 3 {
		return 0, nil, errors.New("msgpackUnpackArrayHeader: too few bytes")
	}
	if b[0] != 0xdc {
		return 0, nil, errors.New("msgpackUnpackArrayHeader: type err")
	}
	return int(uint16(b[1]) + uint16(b[2])), b[3:], nil
}

// for test only
func msgpackUnpackString(b []byte) (string, []byte, error) {
	if len(b) < 3 {
		return "", nil, errors.New("msgpackUnpackString: too few bytes")
	}
	if b[0] != 0xda {
		return "", nil, errors.New("msgpackUnpackString: type err")
	}
	n := int(uint16(b[1]) + uint16(b[2]))
	b = b[3:]
	if len(b) < n {
		return "", nil, errors.New("msgpackUnpackString: too few bytes")
	}
	return string(b[:n]), b[n:], nil
}

// for test only
func msgpackUnpackStringArray(b []byte) ([]string, []byte, error) {
	var n int
	var err error
	n, b, err = msgpackUnpackArrayHeader(b)
	if err != nil {
		return nil, nil, err
	}
	ret := make([]string, 0, n)
	for i := 0; i < n; i++ {
		var s string
		s, b, err = msgpackUnpackString(b)
		if err != nil {
			return nil, nil, err
		}
		ret = append(ret, s)
	}
	return ret, b, nil
}

// for test only
func msgpackUnpack2DStringArray(b []byte) ([][]string, []byte, error) {
	var n int
	var err error
	n, b, err = msgpackUnpackArrayHeader(b)
	if err != nil {
		return nil, nil, err
	}
	ret := make([][]string, 0, n)
	for i := 0; i < n; i++ {
		var ss []string
		ss, b, err = msgpackUnpackStringArray(b)
		if err != nil {
			return nil, nil, err
		}
		ret = append(ret, ss)
	}
	return ret, b, nil
}
