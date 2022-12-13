package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/yogeshlonkar/go-worker-pool"
)

type Job struct {
	str string
}

type Result struct {
	dur time.Duration
	str string
}

type customKey1 struct{}

type customKey2 struct{}

func main() {
	// set context for workers
	ctx := context.WithValue(context.Background(), customKey1{}, "first-value")
	ctx = context.WithValue(ctx, customKey2{}, 2)
	// create a retry able pool with Job and Result type, maxRetry set to 10
	p, err := pool.NewPool(ctx, pool.NewConfig(5, 100, 100, 20, true, worker))
	for i := 0; i < 20; i++ {
		// send a job to pool
		p.SendJobs(Job{fmt.Sprintf("job %d", i)})
	}
	if err != nil {
		panic(err)
	}
	// close jobs, process result from workers
	for result := range p.Close() {
		fmt.Println(result.str, "in", result.dur)
	}
}

// worker with job type Job struct and result type Result
func worker(ctx context.Context, job Job) (Result, error) {
	// use context to get additional parameters/ variables for worker
	someValue1 := ctx.Value(customKey1{}).(string)
	someValue2 := ctx.Value(customKey2{}).(int)
	n := rand.Intn(10-1) + 1
	dur := time.Duration(n) * time.Second
	time.Sleep(dur)
	if n%3 == 0 {
		panic(fmt.Errorf("random panic %d", someValue2))
	} else {
		return Result{
			dur,
			fmt.Sprintf(job.str + " done " + someValue1),
		}, nil
	}
}
