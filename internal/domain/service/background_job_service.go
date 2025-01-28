package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/mburaksoran/insider-case/internal/app/service"
	"github.com/mburaksoran/insider-case/internal/domain/models"
	"github.com/mburaksoran/insider-case/internal/domain/repository"
	"github.com/mburaksoran/insider-case/internal/shared/sqlc_db"
	"log"
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
			log.Printf("Error while creating Job: %v", err)
			return nil, err
		}
		return result, nil
	})

	return err
}

func (s *backgroundJobService) GetActiveJobs(ctx context.Context) ([]*models.BackgroundJob, error) {
	result, err := s.backgroundJobRepository.WithTransaction(ctx, func(q *sqlc_db.Queries) (interface{}, error) {
		jobs, err := s.backgroundJobRepository.GetActiveJobs(ctx, q)
		if err != nil {
			log.Printf("Error while getting active jobs: %v", err)
			return nil, err
		}
		if len(jobs) < 1 {
			return nil, nil
		}
		for _, job := range jobs {
			updateErr := s.backgroundJobRepository.UpdateJobStatus(ctx, q, job.ID, "processing") //TODO create new const for status
			if updateErr != nil {
				return nil, errors.New(fmt.Sprintf("errors while updating job status JobId: %s error: %s", job.ID, err.Error()))
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
			log.Printf("Error while status Job: %v", err)
			return nil, err
		}
		err = s.backgroundJobRepository.UpdateJobLastTriggeredTime(ctx, queries, id)
		if err != nil {
			log.Printf("Error while updating last triggered at Job: %v", err)
			return nil, err
		}
		return nil, nil
	})

	return err
}
