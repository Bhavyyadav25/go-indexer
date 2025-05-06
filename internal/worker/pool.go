package worker

import "sync"

type Pool struct {
	Jobs chan string
	wg   sync.WaitGroup
}

func NewPool(numWorkers int, task func(string) error) *Pool {
	pool := &Pool{
		Jobs: make(chan string),
	}

	pool.wg.Add(numWorkers)
	for i := 0; i < numWorkers; i++ {
		go func() {
			defer pool.wg.Done()
			for job := range pool.Jobs {
				task(job)
			}
		}()
	}
	return pool
}

func (p *Pool) Wait() {
	close(p.Jobs)
	p.wg.Wait()
}
