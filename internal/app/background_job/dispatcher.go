package backgroundjob

import (
	"github.com/mburaksoran/insider-case/internal/app/background_job/tasks"
	"github.com/mburaksoran/insider-case/internal/app/config"
	"github.com/mburaksoran/insider-case/internal/infra/workerpool"
)

type Dispatcher struct {
	ErrorChan  chan error
	Cfg        *config.AppConfig
	WorkerPool workerpool.Pool
}

func NewDispatcher(errorChan chan error, cfg *config.AppConfig, pool workerpool.Pool) *Dispatcher {
	return &Dispatcher{
		ErrorChan:  errorChan,
		Cfg:        cfg,
		WorkerPool: pool,
	}
}

func (d *Dispatcher) DispatchJob(task *tasks.JobTask) {
	d.WorkerPool.AddWork(task)
}
