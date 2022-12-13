package pool

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"testing"
	"time"
)

func TestPool(t *testing.T) {
	start := time.Now()
	worker := func(ctx context.Context, job int) (string, error) {
		time.Sleep(time.Duration(job) * 100 * time.Millisecond)
		return fmt.Sprintf("%d processed", job), nil
	}
	p, err := NewPool(context.Background(), DefaultConfig(5, worker))
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	jobs := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	p.SendJobs(jobs...)
	count := 0
	for result := range p.Close() {
		if !strings.HasSuffix(result, " processed") {
			t.Errorf("expected job '%s' to be processed", result)
		}
		count++
	}
	if count != 10 {
		t.Errorf("missing results expected to be 10, got %d", count)
	}
	if time.Since(start) > 2*time.Second {
		t.Errorf("expected completion time to be 2*time.Second, got %s", time.Since(start))
	}
	if len(p.Errors()) != 0 {
		t.Errorf("expected errors be 0, got %d", len(p.Errors()))
	}
}

func TestPoolErrors(t *testing.T) {
	start := time.Now()
	worker := func(ctx context.Context, job int) (string, error) {
		if job%2 == 0 {
			return "", errors.New("some-error")
		}
		time.Sleep(time.Duration(job) * 100 * time.Millisecond)
		return fmt.Sprintf("%d processed", job), nil
	}
	p, err := NewPool(context.Background(), DefaultConfig(5, worker))
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	jobs := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	p.SendJobs(jobs...)
	count := 0
	for result := range p.Close() {
		if !strings.HasSuffix(result, " processed") {
			t.Errorf("expected job '%s' to be processed", result)
		}
		count++
	}
	if count != 5 {
		t.Errorf("missing results expected to be 5, got %d", count)
	}
	if time.Since(start) > 2*time.Second {
		t.Errorf("expected completion time to be 2*time.Second, got %s", time.Since(start))
	}
	if len(p.Errors()) != 5 {
		t.Errorf("expected errors be 5, got %d", len(p.Errors()))
	}
}

func TestPoolCreationError(t *testing.T) {
	worker := func(ctx context.Context, job int) (string, error) { return "", nil }
	tests := []struct {
		size, jLimit, rLimit, mRetry int
		panic                        bool
		worker                       func(ctx context.Context, job int) (string, error)
		expected                     error
	}{
		{0, 0, 0, 0, false, nil, errors.New("expected pool size to be more than 0")},
		{10, 0, 0, 0, false, nil, errors.New("expected JobQueueLimit to be than 0")},
		{10, 100, 0, 0, false, nil, errors.New("expected ResultQueueLimit to be than 0")},
		{10, 100, 100, 0, false, nil, errors.New("expected worker func to be not nil")},
		{10, 100, 100, 0, false, worker, nil},
	}
	for index, test := range tests {
		t.Run(fmt.Sprintf("ValidationError%d", index), func(t *testing.T) {
			_, err := NewPool(context.Background(), NewConfig(test.size, test.jLimit, test.rLimit, test.mRetry, test.panic, test.worker))
			if err == nil {
				if test.expected != nil {
					t.Errorf("expected error '%s', got nil", test.expected.Error())
				}
			} else {
				if test.expected == nil {
					t.Errorf("expected nil error, got '%s'", err.Error())
				} else if test.expected.Error() != err.Error() {
					t.Errorf("expected error to be '%s', got '%s'", test.expected.Error(), err.Error())
				}
			}

		})
	}
}

func BenchmarkPool(b *testing.B) {
	tests := []struct{ jobs, workers int }{
		{10, 5},
		{10, 10},
		{10, 20},
		{20, 20},
	}
	worker := func(ctx context.Context, job int) (string, error) {
		time.Sleep(time.Duration((rand.Intn(job+1))*100) * time.Millisecond)
		return fmt.Sprintf("%d processed", job), nil
	}
	for _, test := range tests {
		b.Run(fmt.Sprintf("jobs_%d_workers_%d", test.jobs, test.workers), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				p, err := NewPool(context.Background(), DefaultConfig(test.workers, worker))
				if err != nil {
					b.Errorf("expected nil error, got %v", err)
				}
				for job := 0; job < test.jobs; job++ {
					p.SendJobs(job)
				}
				<-p.Close()
			}
		})
	}
}
