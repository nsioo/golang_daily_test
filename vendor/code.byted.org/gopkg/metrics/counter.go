package metrics

import (
	"sync/atomic"
)

// Counter is the interface that provides an optimized version of EmitCounter with fixed tags
// the Incr method
type Counter interface {
	Incr(n int64) int64
	Count() int64
}

type counterImpl struct {
	n, m int64
	e    metricEntity
	c    *metricCache
}

func (p *counterImpl) Incr(n int64) int64 {
	return atomic.AddInt64(&p.n, n)
}

func (p *counterImpl) Count() int64 {
	return atomic.LoadInt64(&p.n)
}

func (p *counterImpl) Flush() error {
	n := p.Count()
	m := atomic.SwapInt64(&p.m, n)
	diff := n - m
	if diff > 0 {
		e := p.e
		e.v = float64(diff)
		return p.c.Send(e)
	}
	return nil
}
