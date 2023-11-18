package future

import (
	"fmt"
	"sync"
	"time"
)

type Completer chan interface{}

func NewCompleteFuture(completer Completer, timeout time.Duration) *CompleteFuture {
	f := &CompleteFuture{complete: completer}
	f.wg.Add(1)
	var wg sync.WaitGroup
	wg.Add(1)
	go listenForResult(f, completer, timeout, &wg)
	wg.Wait()
	return f
}

type CompleteFuture struct {
	complete  Completer
	triggered bool
	result    interface{}
	err       error
	lock      sync.Mutex
	wg        sync.WaitGroup
}

func (f *CompleteFuture) GetFuture() (interface{}, error) {
	f.lock.Lock()
	if f.triggered {
		f.lock.Unlock()
		return f.result, f.err
	}

	f.lock.Unlock()

	f.wg.Wait()
	return f.result, f.err

}

func (f *CompleteFuture) getComplete() Completer {
	return f.complete

}

func (f *CompleteFuture) IsCompleted() bool {
	f.lock.Lock()
	IsCompleted := f.triggered
	f.lock.Unlock()
	return IsCompleted
}

func (f *CompleteFuture) setResult(result interface{}, err error) {
	f.lock.Lock()
	f.triggered = true
	f.result = result
	f.err = err
	f.lock.Unlock()
	f.wg.Done()
}

func listenForResult(f *CompleteFuture, ch Completer, timeout time.Duration, wg *sync.WaitGroup) {
	wg.Done()
	t := time.NewTimer(timeout)

	select {
	case item := <-ch:
		f.setResult(item, nil)
		t.Stop()
	case <-t.C:
		f.setResult(nil, fmt.Errorf(`tomout after %f sencond`, timeout.Seconds()))
	}
}
