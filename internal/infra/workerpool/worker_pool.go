package workerpool

import (
	"errors"
	"go.uber.org/zap"
	"sync"
)

type Pool interface {
	Start()
	Stop()
	AddWork(Task)
}

type Task interface {
	Execute() error
	OnFailure(err error)
}

type WorkerPool struct {
	workerCount int
	lgr         *zap.SugaredLogger
	tasks       chan Task
	start       sync.Once
	stop        sync.Once
	quit        chan struct{}
}

func NewWorkerPool(workerCount int, channelSize int, lgr *zap.SugaredLogger) (Pool, error) {
	if workerCount <= 0 {
		return nil, errors.New("cannot create a pool with less than 0 worker")
	}
	if channelSize < 0 {
		return nil, errors.New("cannot create a pool with less than 0 channel size")
	}

	newTasks := make(chan Task, channelSize)
	return &WorkerPool{
		lgr:         lgr,
		workerCount: workerCount,
		tasks:       newTasks,
		start:       sync.Once{},
		stop:        sync.Once{},
		quit:        make(chan struct{}),
	}, nil
}

func (p *WorkerPool) Start() {
	p.lgr.Info("worker pool started")
	p.start.Do(func() {
		p.startWorkers()
	})
}

func (p *WorkerPool) Stop() {
	p.stop.Do(func() {
		close(p.quit)
	})
}

func (p *WorkerPool) AddWork(t Task) {
	select {
	case p.tasks <- t:
	case <-p.quit:
	}
}

func (p *WorkerPool) startWorkers() {
	for i := 0; i < p.workerCount; i++ {
		go func(workerNum int) {
			for {
				select {
				case <-p.quit:
					return
				case task, ok := <-p.tasks:
					if !ok {
						return
					}
					if err := task.Execute(); err != nil {
						task.OnFailure(err)
					}
				}
			}
		}(i)
	}
}
