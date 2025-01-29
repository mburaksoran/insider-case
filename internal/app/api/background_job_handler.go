package api

import (
	"fmt"
	"github.com/mburaksoran/insider-case/internal/app/service"
	"github.com/mburaksoran/insider-case/internal/domain/models"
	"github.com/mburaksoran/insider-case/internal/domain/types"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/gofiber/fiber/v2"
)

type backgroundJobHandler struct {
	BackgroundService service.BackgroundJobServiceInterface
	logger            *zap.SugaredLogger
}

type BackgroundJobHandlerInterface interface {
	UpdateJobHandler(c *fiber.Ctx) error
	GetJobsHandler(c *fiber.Ctx) error
	StopJobsHandler(c *fiber.Ctx) error
	StartJobsHandler(c *fiber.Ctx) error
}

func NewBackgroundJobHandler(bgService service.BackgroundJobServiceInterface, lgr *zap.SugaredLogger) BackgroundJobHandlerInterface {
	return &backgroundJobHandler{
		BackgroundService: bgService,
		logger:            lgr,
	}
}

func (bh *backgroundJobHandler) UpdateJobHandler(c *fiber.Ctx) error {
	var bj *models.BackgroundJob
	err := c.BodyParser(&bj)
	if err != nil {
		bh.logger.Errorf("[UpdateJobsHandler]  Error while parsing request body, err: %s", err.Error())
		return errors.New(fmt.Sprintf("UpdateJobsHandler  Error while parsing request body, err: %s", err))
	}

	err = bh.BackgroundService.UpdateJob(c.Context(), bj.ID, bj.Status)
	if err != nil {
		bh.logger.Errorf("[UpdateJobsHandler] Error while updating job, err: %s", err.Error())
		return errors.New(fmt.Sprintf("[UpdateJobsHandler] Error while updating job, err: %s", err.Error()))
	}

	return c.SendStatus(fiber.StatusOK)
}

func (bh *backgroundJobHandler) GetJobsHandler(c *fiber.Ctx) error {
	jobs, err := bh.BackgroundService.GetJobs(c.Context())
	if err != nil {
		bh.logger.Errorf("[GetJobsHandler] -  Error while listing jobs, err: %s", err.Error())
		return errors.New(fmt.Sprintf("UpdateJobsHandler Error while listing jobs, err: %s", err.Error()))
	}
	return c.JSON(jobs)
}

func (bh *backgroundJobHandler) StartJobsHandler(c *fiber.Ctx) error {
	err := bh.BackgroundService.UpdateAllJobsStatus(c.Context(), types.BackgroundJobStatusActive)
	if err != nil {
		bh.logger.Errorf("[UpdateAllJobsStatus] Error while updating job status to active err: %s", err.Error())
		return errors.New(fmt.Sprintf("[UpdateAllJobsStatus] Error while updating job status to active , err: %s", err.Error()))
	}
	return c.SendStatus(fiber.StatusOK)
}

func (bh *backgroundJobHandler) StopJobsHandler(c *fiber.Ctx) error {
	err := bh.BackgroundService.UpdateAllJobsStatus(c.Context(), types.BackgroundJobStatusPassive)
	if err != nil {
		bh.logger.Error("[UpdateAllJobsStatus] Error while updating job status to deactivate err: %s", err.Error())
		return errors.New(fmt.Sprintf("[UpdateAllJobsStatus] Error while updating job status to deactivate , err: %s", err.Error()))
	}
	return c.SendStatus(fiber.StatusOK)
}
