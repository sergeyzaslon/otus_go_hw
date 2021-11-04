package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	var errCount uint64
	var wg sync.WaitGroup
	ch := make(chan Task)
	defer close(ch)
	quitCh := make(chan bool)
	defer close(quitCh)

	tasksLen := len(tasks)

	// consumers
	for i := 0; i < n && i < tasksLen; i++ {
		wg.Add(1)
		go func(ch chan Task, quitCh chan bool, i int) {
			for {
				select {
				case task := <-ch:
					{
						if task() != nil {
							atomic.AddUint64(&errCount, 1)
						}
					}
				case <-quitCh:
					wg.Done()
					return
				}
			}
		}(ch, quitCh, i)
	}

	// producer
	for _, task := range tasks {
		if errCount < uint64(m) {
			ch <- task
		} else {
			break
		}
	}

	for i := 0; i < n && i < tasksLen; i++ {
		quitCh <- true
	}

	wg.Wait()

	if m > 0 && errCount >= uint64(m) {
		return ErrErrorsLimitExceeded
	}
	return nil
}
