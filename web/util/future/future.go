package util

import "time"

type Future[T any] struct {
	resultChan chan *T
	errChan    chan error
	timer      *time.Timer
}

func NewFuture[T any]() *Future[T] {
	return &Future[T]{
		resultChan: make(chan *T),
		errChan:    make(chan error),
	}
}

func (f *Future[T]) Result() (*T, error) {
	defer func() {
		close(f.resultChan)
		close(f.errChan)
	}()
	select {
	case res := <-f.resultChan:
		return res, nil
	case err := <-f.errChan:
		return nil, err
	}
}

func (f *Future[T]) Finish(res *T) {
	f.resultChan <- res
}
