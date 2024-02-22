package future

import (
	"errors"
	"time"
)

type Future[T any] struct {
	resultChan chan *T
	errChan    chan error
	timer      *time.Timer
}

var ErrTimeout = errors.New("future timeout")

func NewFuture[T any]() *Future[T] {
	return &Future[T]{
		resultChan: make(chan *T),
		errChan:    make(chan error),
	}
}

func WithTimeout[T any](timeout time.Duration) *Future[T] {
	return &Future[T]{
		resultChan: make(chan *T),
		errChan:    make(chan error),
		timer:      time.NewTimer(timeout),
	}
}

func (f *Future[T]) Result() (*T, error) {
	defer func() {
		close(f.resultChan)
		close(f.errChan)
		f.timer.Stop()
	}()
	if f.timer == nil {
		select {
		case res := <-f.resultChan:
			return res, nil
		case err := <-f.errChan:
			return nil, err
		}
	} else {
		select {
		case res := <-f.resultChan:
			return res, nil
		case err := <-f.errChan:
			return nil, err
		case <-f.timer.C:
			return nil, ErrTimeout
		}
	}
}

func (f *Future[T]) Success(res *T) {
	f.resultChan <- res
}

func (f *Future[T]) Fail(err error) {
	f.errChan <- err
}
