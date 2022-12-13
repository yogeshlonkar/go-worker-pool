package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/yogeshlonkar/go-worker-pool"
)

func main() {
	// create a pool with string job and string result type
	p, err := pool.NewPool(context.Background(), pool.DefaultConfig(5, worker))
	if err != nil {
		panic(err)
	}
	for i := 0; i < 20; i++ {
		// send a job to pool
		p.SendJobs(fmt.Sprintf("job %d", i))
	}
	// close jobs for workers to exit after finishing jobs and get results channel
	results := p.Close()
	// process result from workers, will wait until all workers exit
	for result := range results {
		fmt.Println(result)
	}
}

// worker that processes string job and returns string result
func worker(ctx context.Context, job string) (string, error) {
	n := rand.Intn(3-1) + 1
	time.Sleep(time.Duration(n) * time.Second)
	return fmt.Sprintf(job + " done"), nil
}
