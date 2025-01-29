package backgroundjob

import (
	"context"
	"database/sql"
	"github.com/mburaksoran/insider-case/internal/app/service"
	"github.com/mburaksoran/insider-case/internal/domain/models"
	"go.uber.org/zap"
	"time"
)

type Fetcher struct {
	DB                   *sql.DB
	BackgroundJobService service.BackgroundJobServiceInterface
	JobChan              chan *models.BackgroundJob
	ErrorChan            chan error
	logger               *zap.SugaredLogger
}

func NewFetcher(db *sql.DB, bgService service.BackgroundJobServiceInterface, jobChan chan *models.BackgroundJob, errorChan chan error, lgr *zap.SugaredLogger) *Fetcher {
	return &Fetcher{
		DB:                   db,
		BackgroundJobService: bgService,
		JobChan:              jobChan,
		ErrorChan:            errorChan,
		logger:               lgr,
	}
}

func (f *Fetcher) Start() {
	ticker := time.NewTicker(10 * time.Second)
	f.logger.Info("Checking Active Jobs")
	defer ticker.Stop()
	for range ticker.C {
		f.FetchJobs()
	}
}

func (f *Fetcher) FetchJobs() {
	ctx := context.Background()
	dueJobs, err := f.BackgroundJobService.GetActiveJobsForBackgroundService(ctx)
	if err != nil {
		f.ErrorChan <- err
		return
	}

	if len(dueJobs) != 0 {
		f.logger.Infof("Active Jobs Found Count: %v", len(dueJobs))
		for _, job := range dueJobs {
			f.JobChan <- job
		}
	}

}
