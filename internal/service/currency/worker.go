package currency

import (
	"context"
	"sync"
)

type doTaskFunc func() interface{}

// Do your tasks parallel.
// <-chan stream results each of your funcs.
func (s *Currency) taskResultStream(ctx context.Context, tasks ...doTaskFunc) <-chan interface{} {
	out := make(chan interface{}, len(tasks))

	var wgClose sync.WaitGroup
	wgClose.Add(len(tasks))
	go func() {
		wgClose.Wait()
		close(out)
	}()

	for _, f := range tasks {
		go func() {
			defer wgClose.Done()
			result := f()
			select {
			case <-ctx.Done():
			case out <- result:
			}
		}()
	}

	return out
}
