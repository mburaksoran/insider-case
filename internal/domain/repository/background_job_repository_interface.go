package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/mburaksoran/insider-case/internal/domain/models"
	"github.com/mburaksoran/insider-case/internal/shared/sqlc_db"
)

type BackgroundJobRepositoryInterface interface {
	CreateJob(ctx context.Context, queries *sqlc_db.Queries, job models.BackgroundJob) (bool, error)
	GetActiveJobsForBackgroundService(ctx context.Context, queries *sqlc_db.Queries) ([]*models.BackgroundJob, error)
	UpdateJobLastTriggeredTime(ctx context.Context, queries *sqlc_db.Queries, id uuid.UUID) error
	UpdateJobStatus(ctx context.Context, queries *sqlc_db.Queries, id uuid.UUID, status string) error
	WithTransaction(ctx context.Context, fn func(*sqlc_db.Queries) (interface{}, error)) (interface{}, error)
	WithoutTransaction(ctx context.Context, f func(*sqlc_db.Queries) (interface{}, error)) (interface{}, error)
	GetJobs(ctx context.Context, queries *sqlc_db.Queries) ([]*models.BackgroundJob, error)
	UpdateAllJobsStatus(ctx context.Context, queries *sqlc_db.Queries, status string) error
}
