package parallel

import (
	"sync"
)

func Parallel(fn func(int) any, times int) []any {
	var wg sync.WaitGroup
	var results = make([]any, times)
	for i := 0; i < times; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			results[index] = fn(index)
		}(i)
	}

	wg.Wait()
	return results
}
