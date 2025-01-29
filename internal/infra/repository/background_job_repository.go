package repository

import (
	"context"
	"database/sql"
	"github.com/pkg/errors"
	"time"

	"github.com/google/uuid"
	"github.com/mburaksoran/insider-case/internal/domain/models"
	"github.com/mburaksoran/insider-case/internal/domain/repository"
	"github.com/mburaksoran/insider-case/internal/shared/sqlc_db"
)

type backgroundJobRepository struct {
	db *sql.DB
}

func NewBackgroundJobRepository(db *sql.DB) repository.BackgroundJobRepositoryInterface {
	return &backgroundJobRepository{db: db}
}

func (r *backgroundJobRepository) CreateJob(ctx context.Context, queries *sqlc_db.Queries, job models.BackgroundJob) (bool, error) {

	_, err := queries.CreateJob(ctx, sqlc_db.CreateJobParams{
		Name:          job.Name,
		Handler:       job.Handler,
		Interval:      job.Interval,
		Status:        job.Status,
		LastTriggered: sql.NullTime{},
	})

	if err != nil {
		return false, errors.Wrapf(err, "[CreateJob] - CreateJob Error")
	}
	return true, nil
}

// It should only be used in background jobs because it change the status after pulling jobs
func (r *backgroundJobRepository) GetActiveJobsForBackgroundService(ctx context.Context, queries *sqlc_db.Queries) ([]*models.BackgroundJob, error) {
	list, err := queries.GetDueJobs(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "[GetActiveJobsForBackgroundService] - GetDueJobs Error")
	}
	if len(list) < 1 {
		return nil, nil
	}

	var resultList []*models.BackgroundJob

	for _, job := range list {
		resultList = append(resultList, &models.BackgroundJob{
			ID:            job.ID,
			Name:          job.Name,
			Handler:       job.Handler,
			Interval:      job.Interval,
			Status:        job.Status,
			LastTriggered: job.LastTriggered.Time,
		})
	}

	return resultList, nil
}

func (r *backgroundJobRepository) UpdateJobLastTriggeredTime(ctx context.Context, queries *sqlc_db.Queries, id uuid.UUID) error {

	err := queries.UpdateJobLastTriggered(ctx, sqlc_db.UpdateJobLastTriggeredParams{
		LastTriggered: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
		ID: id,
	})

	if err != nil {
		return errors.Wrapf(err, "[UpdateJobLastTriggeredTime] - UpdateJobLastTriggered Error")
	}

	return nil
}

func (r *backgroundJobRepository) UpdateJobStatus(ctx context.Context, queries *sqlc_db.Queries, id uuid.UUID, status string) error {

	err := queries.UpdateJobStatus(ctx, sqlc_db.UpdateJobStatusParams{
		Status: status,
		ID:     id,
	})

	if err != nil {
		return errors.Wrapf(err, "[UpdateJobStatus] - UpdateJobStatus Error")
	}

	return nil
}

func (r *backgroundJobRepository) UpdateAllJobsStatus(ctx context.Context, queries *sqlc_db.Queries, status string) error {

	err := queries.UpdateAllJobsStatus(ctx, status)
	if err != nil {
		return errors.Wrapf(err, "[UpdateAllJobsStatus] - UpdateAllJobsStatus Error")
	}
	return nil
}

func (r *backgroundJobRepository) GetJobs(ctx context.Context, queries *sqlc_db.Queries) ([]*models.BackgroundJob, error) {

	result, err := queries.GetJobs(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "[GetJobs] - GetJobs Error")
	}
	var jobList []*models.BackgroundJob

	for _, job := range result {
		jobList = append(jobList, &models.BackgroundJob{
			ID:            job.ID,
			Name:          job.Name,
			Handler:       job.Handler,
			Interval:      job.Interval,
			Status:        job.Status,
			LastTriggered: job.LastTriggered.Time,
		})
	}

	return jobList, nil
}

func (r *backgroundJobRepository) WithoutTransaction(ctx context.Context, fn func(*sqlc_db.Queries) (interface{}, error)) (interface{}, error) {
	q := sqlc_db.New(r.db)
	return fn(q)
}

func (r *backgroundJobRepository) WithTransaction(ctx context.Context, fn func(*sqlc_db.Queries) (interface{}, error)) (interface{}, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, errors.Wrapf(err, "[WithTransaction] - BeginTx Error")
	}
	q := sqlc_db.New(tx)
	res, err := fn(q)
	if err != nil {
		tx.Rollback()
		return nil, errors.Wrapf(err, "[WithTransaction] - fn(q) Error")
	}
	return res, tx.Commit()
}
