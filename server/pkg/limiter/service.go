package limiter

import (
	"sync"
	"time"
)

type Limiter interface {
	Set(addr string) *BlockData
	Unblock(addr string)
	StartCleaner()
}

type Data struct {
	m    *sync.Mutex
	data map[string]*BlockData
}

type BlockData struct {
	count   int
	Blocked bool
}

func New() Limiter {
	return &Data{
		m:    &sync.Mutex{},
		data: make(map[string]*BlockData, 10),
	}
}

func (l *Data) Set(addr string) *BlockData {
	l.m.Lock()
	defer l.m.Unlock()

	b, ok := l.data[addr]
	if !ok {

		b = &BlockData{}
		l.data[addr] = b
	}

	if b.count > 5 {
		b.Blocked = true
	}

	b.count += 1
	return b
}

func (l *Data) Unblock(addr string) {
	l.m.Lock()
	defer l.m.Unlock()

	l.data[addr].Blocked = false
	l.data[addr].count = 0
}

func (l *Data) StartCleaner() {
	go func() {
		ticker := time.NewTicker(time.Second * 5)
		for range ticker.C {
			l.m.Lock()

			for _, v := range l.data {
				v.count = 0
			}

			l.m.Unlock()
		}
	}()
}
