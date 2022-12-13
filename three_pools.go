package pool

import (
	"context"
	"sync"
)

type threeStagePool[J, R1, R2, R3 any] struct {
	p1 *singleStagePool[J, R1]
	p2 *singleStagePool[R1, R2]
	p3 *singleStagePool[R2, R3]
}

// NewThreeStagePool creates new instance of three chained worker pools and starts workers
func NewThreeStagePool[J, R1, R2, R3 any](ctx context.Context, config1 *Config[J, R1], config2 *Config[R1, R2], config3 *Config[R2, R3]) (Pool[J, R3], error) {
	p := &threeStagePool[J, R1, R2, R3]{
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
		p3: &singleStagePool[R2, R3]{
			Config:  config3,
			running: config3.Size,
			mutex:   sync.Mutex{},
		},
	}
	if err := p.validate(); err != nil {
		return nil, err
	}
	p.startPool(ctx)
	return p, nil
}

func (p *threeStagePool[J, R1, R2, R3]) validate() error {
	if err := p.p1.validate(); err != nil {
		return err
	}
	if err := p.p2.validate(); err != nil {
		return err
	}
	return p.p3.validate()
}

func (p *threeStagePool[J, R1, R2, R3]) startPool(ctx context.Context) {
	p.p1.startPool(ctx, make(chan J, p.p1.JobQueueLimit), make(chan R1, p.p1.ResultQueueLimit))
	p.p2.startPool(ctx, p.p1.results, make(chan R2, p.p2.ResultQueueLimit))
	p.p3.startPool(ctx, p.p2.results, make(chan R3, p.p3.ResultQueueLimit))
}

// SendJobs to job que for first worker pool
func (p *threeStagePool[J, R1, R2, R3]) SendJobs(jobs ...J) {
	for _, job := range jobs {
		p.p1.jobs <- job
	}
}

// Close closes job que and returns results channel for 3rd worker pool
func (p *threeStagePool[J, R1, R2, R3]) Close() <-chan R3 {
	p.p1.Close()
	return p.p3.results
}

func (p *threeStagePool[J, R1, R2, R3]) Errors() []JobError {
	errors := append(p.p1.Errors(), p.p2.Errors()...)
	return append(errors, p.p3.Errors()...)
}
