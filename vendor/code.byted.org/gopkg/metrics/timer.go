package metrics

import (
	"sync"
	"sync/atomic"
	"time"
)

// Timer is the interface that provides an optimized version of EmitTimer with fixed tags
// it saves all added values with compacted memory, and flushes all data to tsdb
type Timer interface {
	AddDuration(time.Duration) // it's same as Add(float64(t)), the unit of tsdb is nanosecond
	Add(float64)
	Count() int64 // counter of Add && AddDuration
}

const (
	float64sShards = 16
)

type float64s struct {
	sync.Mutex
	ff []float64
}

type timerImpl struct {
	n, m   int64
	shards [float64sShards]float64s

	e metricEntity
	c *metricCache
}

func (p *timerImpl) AddDuration(t time.Duration) {
	p.Add(float64(t))
}

func (p *timerImpl) Count() int64 {
	return atomic.LoadInt64(&p.n)
}

func (p *timerImpl) Add(v float64) {
	idx := atomic.AddInt64(&p.n, 1) % float64sShards
	s := &p.shards[idx]
	s.Lock()
	if s.ff == nil {
		s.ff = make([]float64, 0, maxPendingSize)
	}
	s.ff = append(s.ff, v)
	var flushff []float64 // flush without lock
	if len(s.ff) >= maxPendingSize {
		flushff = append(make([]float64, 0, maxPendingSize), s.ff...)
		s.ff = s.ff[:0]
	}
	s.Unlock()

	if len(flushff) == 0 {
		return
	}

	// flush data if shard is full.
	if err := p.sendff(flushff); err != nil {
		logfunc("[ERROR] gopkg/metrics: send err: %s", err)
	}
}

func (p *timerImpl) sendff(ff []float64) error {
	// optimize for small slice try to cache it
	if len(ff) < maxPendingSize/10 {
		for _, v := range ff {
			e := p.e
			e.v = v
			if err := p.c.Send(e); err != nil {
				return err
			}
		}
		return nil
	}
	// no need to use `metricCache.Send` one by one, and it's quite slow
	ee := mesPoolGet()
	for _, v := range ff {
		e := p.e
		e.v = v
		*ee = append(*ee, e)
	}
	return p.c.send(ee)
}

func (p *timerImpl) Flush() error {
	n, m := atomic.LoadInt64(&p.n), atomic.LoadInt64(&p.m)
	if n == m { // optimize for empty timers
		return nil
	}
	atomic.StoreInt64(&p.m, n)

	ff := make([]float64, 0, 2*maxPendingSize)
	for i := range p.shards {
		p.shards[i].Lock()
		ff = append(ff, p.shards[i].ff...)
		p.shards[i].ff = p.shards[i].ff[:0]
		p.shards[i].Unlock()

		if len(ff) >= maxPendingSize {
			if err := p.sendff(ff); err != nil {
				return err
			}
			ff = ff[:0]
		}
	}
	return p.sendff(ff)
}
