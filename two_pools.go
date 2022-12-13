package pool

import (
	"context"
	"sync"
)

type twoStagePool[J, R1, R2 any] struct {
	p1 *singleStagePool[J, R1]
	p2 *singleStagePool[R1, R2]
}

// NewTwoStagePool creates new instance of two chained worker pools and starts workers
func NewTwoStagePool[J, R1, R2 any](ctx context.Context, config1 *Config[J, R1], config2 *Config[R1, R2]) (Pool[J, R2], error) {
	p := &twoStagePool[J, R1, R2]{
		p1: &singleStagePool[J, R1]{
			Config:  config1,
			running: config1.Size,
			mutex:   sync.Mutex{},
		},
		p2: &singleStagePool[R1, R2]{
			Config:  config2,
			running: config2.Size,
			mutex:   sync.Mutex{},
		},
	}
	if err := p.validate(); err != nil {
		return nil, err
	}
	p.startPool(ctx)
	return p, nil
}

func (p *twoStagePool[J, R1, R2]) validate() error {
	if err := p.p1.validate(); err != nil {
		return err
	}
	return p.p2.validate()
}

func (p *twoStagePool[J, R1, R2]) startPool(ctx context.Context) {
	p.p1.startPool(ctx, make(chan J, p.p1.JobQueueLimit), make(chan R1, p.p1.ResultQueueLimit))
	p.p2.startPool(ctx, p.p1.results, make(chan R2, p.p2.ResultQueueLimit))
}

// SendJobs to job que for first worker pool
func (p *twoStagePool[J, R1, R2]) SendJobs(jobs ...J) {
	for _, job := range jobs {
		p.p1.jobs <- job
	}
}

// Close closes job que and returns results channel for 2nd worker pool
func (p *twoStagePool[J, R1, R2]) Close() <-chan R2 {
	p.p1.Close()
	return p.p2.results
}

func (p *twoStagePool[J, R1, R2]) Errors() []JobError {
	return append(p.p1.Errors(), p.p2.Errors()...)
}
