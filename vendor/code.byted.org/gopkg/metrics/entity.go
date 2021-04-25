package metrics

import "unsafe"

type metricEntity struct {
	mt     metricsType
	prefix string
	name   string
	v      float64
	ts     int64

	tt *cachedTags
}

type metricEntityKey struct { // for deduplicating
	mt     metricsType
	prefix string
	name   string
	tt     uintptr // it's ok, since we always reuse *cachedTags within a tagCache
}

func (m *metricEntity) Key() metricEntityKey {
	return metricEntityKey{
		mt:     m.mt,
		prefix: m.prefix,
		name:   m.name,
		tt:     uintptr(unsafe.Pointer(m.tt)),
	}
}

func (m *metricEntity) IsCounter() bool {
	return m.mt == metricsTypeCounter || m.mt == metricsTypeRateCounter || m.mt == metricsTypeMeter
}

func (m *metricEntity) MarshalSize() int {
	// protocol: 6 fields: emit $type $prefix.name  $value $tag ""
	n := 0
	n += msgpackArrayHeaderSize
	n += msgpackStringSize(_emit)
	n += msgpackStringSize(m.mt.String())
	if len(m.prefix) > 0 {
		n += msgpackStringHeaderSize + (len(m.prefix) + 1 + len(m.name))
	} else {
		n += msgpackStringSize(m.name)
	}
	n += msgpackStringHeaderSize + floatStrSize(m.v) // int64 + "." + 5 prec float + str header
	n += msgpackStringHeaderSize + len(m.tt.Bytes())
	if m.ts > 0 {
		n += msgpackStringHeaderSize + int64StrSize(m.ts)
	} else {
		n += msgpackStringHeaderSize + 0
	}
	return n
}

func (m *metricEntity) AppendTo(p []byte) []byte {
	// protocol: 6 fields: emit $type $prefix.name  $value $tag ""
	p = msgpackAppendArrayHeader(p, 6)
	p = msgpackAppendString(p, _emit)
	p = msgpackAppendString(p, m.mt.String())
	if len(m.prefix) > 0 {
		p = msgpackAppendStringHeader(p, uint16(len(m.prefix)+1+len(m.name)))
		p = append(p, m.prefix...)
		p = append(p, '.')
		p = append(p, m.name...)
	} else {
		p = msgpackAppendString(p, m.name)
	}
	p = msgpackAppendStringHeader(p, uint16(floatStrSize(m.v)))
	p = appendFloat64(p, m.v)
	p = msgpackAppendStringHeader(p, uint16(len(m.tt.Bytes())))
	p = append(p, m.tt.Bytes()...)
	if m.ts > 0 {
		p = msgpackAppendStringHeader(p, uint16(int64StrSize(m.ts)))
		p = appendInt64(p, m.ts)
	} else {
		p = msgpackAppendString(p, "")
	}
	return p
}

func (m *metricEntity) MarshalTo(b []byte) {
	p := b[:0]
	p = m.AppendTo(p)
	if len(p) != len(b) {
		panic("buf size err")
	}
}
