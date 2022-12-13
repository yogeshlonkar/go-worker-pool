package main

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/yogeshlonkar/go-worker-pool"
)

func main() {
	config1 := pool.NewConfig(5, 100, 100, 10, false, worker1)
	config2 := pool.NewConfig(5, 100, 100, 0, false, worker2)
	// create a pool with string job and string result type
	p, err := pool.NewTwoStagePool(context.Background(), config1, config2)
	for i := 0; i < 20; i++ {
		// send a job to pool
		p.SendJobs(fmt.Sprintf("job %d", i))
	}
	if err != nil {
		panic(err)
	}
	// close jobs for workers to exit after finishing jobs and get results channel
	results := p.Close()
	// process result from workers, will wait until all workers exit
	for result := range results {
		fmt.Println(result)
	}
	for _, err := range p.Errors() {
		fmt.Printf("%s -> %q\n", err.Job, err.Err)
	}
}

// worker1 that processes string job and returns string result
func worker1(_ context.Context, job string) (string, error) {
	n := rand.Intn(3-1) + 1
	if n%2 == 0 {
		return "", errors.New("failed at worker1")
	}
	time.Sleep(time.Duration(n) * time.Second)
	fmt.Println(job + " processed")
	return fmt.Sprintf(job + " one: processed"), nil
}

// worker2 that processes string job and returns string result
func worker2(_ context.Context, job string) (string, error) {
	n := rand.Intn(3-1) + 1
	time.Sleep(time.Duration(n) * time.Second)
	return fmt.Sprintf(job + " done"), nil
}
