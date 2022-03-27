package locks

import (
	"TestTask-ProfOfWork/server/pkg/tasks"
	"sync"
)

//TODO implement to REDIS

type Service interface {
	Get(key string) (*Lock, bool)
	Add(key string, task *Lock)
	Reset(key string)
	AnswerAttempt(key string)
}

type Locks struct {
	m    *sync.Mutex
	data map[string]*Lock
}

type Lock struct {
	Task          *tasks.Task
	AnswerAttempt int
}

func New() *Locks {
	return &Locks{
		m:    &sync.Mutex{},
		data: make(map[string]*Lock, 100),
	}
}

func (l *Locks) Get(key string) (*Lock, bool) {
	l.m.Lock()
	defer l.m.Unlock()

	val, ok := l.data[key]
	if !ok {
		return nil, false
	}

	return val, true
}

func (l *Locks) Add(key string, task *Lock) {
	l.m.Lock()
	defer l.m.Unlock()

	l.data[key] = task
}

func (l *Locks) Reset(key string) {
	l.m.Lock()
	defer l.m.Unlock()

	delete(l.data, key)
}

func (l *Locks) AnswerAttempt(key string) {
	l.m.Lock()
	defer l.m.Unlock()

	l.data[key].AnswerAttempt += 1
}
