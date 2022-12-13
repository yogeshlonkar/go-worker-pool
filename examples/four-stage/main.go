package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/yogeshlonkar/go-worker-pool"
)

func main() {
	config1 := pool.NewConfig(5, 100, 100, 0, false, worker1)
	config2 := pool.NewConfig(5, 100, 100, 0, false, worker2)
	config3 := pool.NewConfig(5, 100, 100, 0, false, worker3)
	config4 := pool.NewConfig(5, 100, 100, 0, false, worker4)
	// create a pool with string job and string result type
	p, err := pool.NewFourStagePool(context.Background(), config1, config2, config3, config4)
	for i := 0; i < 10; i++ {
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
}

// worker1 that processes string job and returns string result
func worker1(ctx context.Context, job string) (string, error) {
	n := rand.Intn(9-1) + 1
	time.Sleep(time.Duration(n) * time.Second)
	fmt.Println(job + " one: processed")
	return fmt.Sprintf(job + " one: processed"), nil
}

// worker2 that processes string job and returns string result
func worker2(ctx context.Context, job string) (string, error) {
	n := rand.Intn(6-1) + 1
	time.Sleep(time.Duration(n) * time.Second)
	fmt.Println(job + " two: processed")
	return fmt.Sprintf(job + " two: processed"), nil
}

// worker3 that processes string job and returns string result
func worker3(ctx context.Context, job string) (string, error) {
	n := rand.Intn(3-1) + 1
	time.Sleep(time.Duration(n) * time.Second)
	fmt.Println(job + " three: processed")
	return fmt.Sprintf(job + " three: processed"), nil
}

// worker4 that processes string job and returns string result
func worker4(ctx context.Context, job string) (string, error) {
	return fmt.Sprintf(job + " done"), nil
}
