package limiter

import (
	"sync"
	"sync/atomic"
)

type Limiter struct {
	left int64
	ch   chan struct{}
	wg   sync.WaitGroup
}

func New(num int64) *Limiter {
	l := &Limiter{
		left: num,
		ch:   make(chan struct{}, num),
	}
	for i := 0; int64(i) < num; i++ {
		l.ch <- struct{}{}
	}
	return l
}

func (l *Limiter) Allow() {
	<-l.ch
	atomic.AddInt64(&l.left, -1)
	l.wg.Add(1)
}

func (l *Limiter) Done() {
	l.ch <- struct{}{}
	atomic.AddInt64(&l.left, 1)
	l.wg.Done()
}

func (l *Limiter) Left() int64 {
	return l.left
}

func (l *Limiter) Wait() {
	l.wg.Wait()
}
