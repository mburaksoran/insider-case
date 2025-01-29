package tasks

import (
	"context"
	"fmt"
	"go.uber.org/zap"

	"github.com/mburaksoran/insider-case/internal/app/background_job/handler"
	"github.com/mburaksoran/insider-case/internal/domain/models"
)

type JobTask struct {
	Job     *models.BackgroundJob
	Factory *handler.HandlerFactory
	Ctx     context.Context
	ErrorCh chan error
	logger  *zap.SugaredLogger
}

func NewJobTask(job *models.BackgroundJob, factory *handler.HandlerFactory, lgr *zap.SugaredLogger) *JobTask {
	return &JobTask{
		Job:     job,
		Factory: factory,
		logger:  lgr,
	}
}

func (t *JobTask) Execute() error {
	jobHandler, err := t.Factory.GetHandler(t.Job.Handler)
	if err != nil {
		t.logger.Errorf("handler not found: %s", t.Job.Handler)
		return fmt.Errorf("handler not found: %s", t.Job.Handler)
	}
	return jobHandler.Handle(t.Ctx, t.Job.ID)
}

func (t *JobTask) OnFailure(err error) {
	t.ErrorCh <- err
}
