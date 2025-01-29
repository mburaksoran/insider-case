package service

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/mburaksoran/insider-case/internal/domain/types"
	"github.com/pkg/errors"

	"github.com/mburaksoran/insider-case/internal/app/service"
	"github.com/mburaksoran/insider-case/internal/domain/models"
	"github.com/mburaksoran/insider-case/internal/domain/repository"
	"github.com/mburaksoran/insider-case/internal/shared/sqlc_db"
)

type backgroundJobService struct {
	backgroundJobRepository repository.BackgroundJobRepositoryInterface
	httpClient              service.HttpClientInterface
}

func NewBackgroundJobService(bgRepo repository.BackgroundJobRepositoryInterface) service.BackgroundJobServiceInterface {
	return &backgroundJobService{
		backgroundJobRepository: bgRepo,
	}
}

func (s *backgroundJobService) CreateJob(ctx context.Context, job models.BackgroundJob) error {
	_, err := s.backgroundJobRepository.WithoutTransaction(ctx, func(queries *sqlc_db.Queries) (interface{}, error) {
		result, err := s.backgroundJobRepository.CreateJob(ctx, queries, job)
		if err != nil {
			return nil, errors.Wrapf(err, "[GetSendMessage] - GetSendMessages Error")
		}
		return result, nil
	})

	return err
}

func (s *backgroundJobService) GetActiveJobsForBackgroundService(ctx context.Context) ([]*models.BackgroundJob, error) {
	result, err := s.backgroundJobRepository.WithTransaction(ctx, func(q *sqlc_db.Queries) (interface{}, error) {
		jobs, err := s.backgroundJobRepository.GetActiveJobsForBackgroundService(ctx, q)
		if err != nil {
			return nil, errors.Wrapf(err, "[GetActiveJobsForBackgroundService] - GetActiveJobsForBackgroundService Error")
		}
		if len(jobs) < 1 {
			return nil, nil
		}
		for _, job := range jobs {
			updateErr := s.backgroundJobRepository.UpdateJobStatus(ctx, q, job.ID, types.BackgroundJobStatusInProgress)
			if updateErr != nil {
				return nil, errors.Wrapf(err, fmt.Sprintf("[UpdateJobStatus] - UpdateJobStatus Error JobId: %s", job.ID))
			}
		}
		return jobs, nil
	})

	if err != nil {
		return nil, err
	}

	if result != nil {
		jobList := result.([]*models.BackgroundJob)
		return jobList, nil
	}

	return nil, nil
}

func (s *backgroundJobService) ActivateJob(ctx context.Context, id uuid.UUID, status string) error {
	_, err := s.backgroundJobRepository.WithoutTransaction(ctx, func(queries *sqlc_db.Queries) (interface{}, error) {
		err := s.backgroundJobRepository.UpdateJobStatus(ctx, queries, id, status)
		if err != nil {
			return nil, errors.Wrapf(err, "[ActivateJob] - UpdateJobStatus Error")
		}
		err = s.backgroundJobRepository.UpdateJobLastTriggeredTime(ctx, queries, id)
		if err != nil {
			return nil, errors.Wrapf(err, "[ActivateJob] - UpdateJobLastTriggeredTime Error")
		}
		return nil, nil
	})

	return err
}

func (s *backgroundJobService) UpdateJob(ctx context.Context, id uuid.UUID, status string) error {
	_, err := s.backgroundJobRepository.WithoutTransaction(ctx, func(queries *sqlc_db.Queries) (interface{}, error) {
		err := s.backgroundJobRepository.UpdateJobStatus(ctx, queries, id, status)
		if err != nil {

			return nil, errors.Wrapf(err, "[UpdateJob] - UpdateJobStatus Error")
		}
		return nil, nil
	})

	return err
}
func (s *backgroundJobService) UpdateAllJobsStatus(ctx context.Context, status string) error {
	_, err := s.backgroundJobRepository.WithoutTransaction(ctx, func(queries *sqlc_db.Queries) (interface{}, error) {
		err := s.backgroundJobRepository.UpdateAllJobsStatus(ctx, queries, status)
		if err != nil {
			return nil, errors.Wrapf(err, "[UpdateAllJobsStatus] - UpdateAllJobsStatus Error")
		}
		return nil, nil
	})

	return err
}
func (s *backgroundJobService) GetJobs(ctx context.Context) ([]*models.BackgroundJob, error) {
	res, err := s.backgroundJobRepository.WithoutTransaction(ctx, func(queries *sqlc_db.Queries) (interface{}, error) {
		result, err := s.backgroundJobRepository.GetJobs(ctx, queries)
		if err != nil {
			return nil, errors.Wrapf(err, "[GetJobs] - GetJobs Error")
		}
		return result, nil
	})
	var jobList []*models.BackgroundJob
	if res != nil {
		jobList = res.([]*models.BackgroundJob)
		return jobList, err
	}
	return nil, nil
}
