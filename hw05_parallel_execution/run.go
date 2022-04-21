package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

type taskError struct {
	mu     *sync.Mutex
	errors int
}

func (t *taskError) addError() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.errors++
}

func (t *taskError) checkLimit(m int) bool {
	t.mu.Lock()
	defer t.mu.Unlock()

	return t.errors >= m
}

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	t := taskError{
		errors: 0,
		mu:     &sync.Mutex{},
	}

	if t.checkLimit(m) {
		return ErrErrorsLimitExceeded
	}

	var taskCh = make(chan Task, len(tasks))

	for _, task := range tasks {
		taskCh <- task
	}
	close(taskCh)

	wg := sync.WaitGroup{}
	wg.Add(n)

	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()
			for task := range taskCh {
				if t.checkLimit(m) {
					break
				}

				resultTask := task()

				if resultTask != nil {
					t.addError()
				}
			}
		}()
	}

	wg.Wait()

	if t.checkLimit(m) {
		return ErrErrorsLimitExceeded
	}

	return nil
}
