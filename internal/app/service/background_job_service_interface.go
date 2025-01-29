package service

import (
	"context"
	"github.com/google/uuid"
	"github.com/mburaksoran/insider-case/internal/domain/models"
)

type BackgroundJobServiceInterface interface {
	CreateJob(ctx context.Context, job models.BackgroundJob) error
	GetActiveJobsForBackgroundService(ctx context.Context) ([]*models.BackgroundJob, error)
	ActivateJob(ctx context.Context, id uuid.UUID, status string) error
	UpdateJob(ctx context.Context, id uuid.UUID, status string) error
	GetJobs(ctx context.Context) ([]*models.BackgroundJob, error)
	UpdateAllJobsStatus(ctx context.Context, status string) error
}
