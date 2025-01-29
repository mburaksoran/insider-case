package handler

import (
	"context"
	"github.com/google/uuid"
	"github.com/mburaksoran/insider-case/internal/app/service"
	"github.com/mburaksoran/insider-case/internal/domain/models"
	"github.com/mburaksoran/insider-case/internal/domain/types"
	"go.uber.org/zap"
	"time"
)

type MessagePublishHandler struct {
	BackgroundService   service.BackgroundJobServiceInterface
	MessageService      service.MessageServiceInterface
	RedisMessageService service.RedisMessageServiceInterface
	HttpClient          service.HttpClientInterface
	logger              *zap.SugaredLogger
}

func (rc *MessagePublishHandler) Handle(ctx context.Context, values interface{}) error {
	messages, err := rc.MessageService.GetMessageToSend(ctx)
	if err != nil {
		rc.logger.Errorf("[GetMessageToSend] - Error while collecting messages for sending err: %s", err.Error())
		return err
	}
	if len(messages) < 1 {
		return nil
	}
	for _, msg := range messages {
		_, httpErr := rc.HttpClient.PostWithAPIKey(ctx, msg.MapToDto())
		if httpErr != nil {
			rc.logger.Errorf("[PostWithAPIKey] - Error while posting messages to the api err: %s", httpErr.Error())
			err = rc.MessageService.UpdateMessageStatus(ctx, msg.ID, types.MessageStatusPending)
			if err != nil {
				rc.logger.Errorf("[UpdateMessageStatus] - Error while updating message status for rollback err: %s", httpErr.Error())
				return err
			}
			return httpErr
		}
		t := models.MessageReceiveHistory{
			ID:          uuid.New(),
			SendingTime: time.Now(),
		}
		redisErr := rc.RedisMessageService.SetMessageReceiveHistory(ctx, t)
		if redisErr != nil {
			rc.logger.Errorf("[SetMessageReceiveHistory] - Error while setting message receive history to redis err: %s", redisErr.Error())
		}
		err = rc.MessageService.UpdateMessageStatus(ctx, msg.ID, types.MessageStatusSent)
		if err != nil {
			rc.logger.Errorf("[UpdateMessageStatus] - Error hile updating message status to sended err: %s", err.Error())
			return err
		}
	}
	err = rc.BackgroundService.ActivateJob(ctx, values.(uuid.UUID), "active")
	if err != nil {
		rc.logger.Errorf("[ActivateJob] - Error while activating job err: %s", err.Error())
		return err
	}
	return nil
}

func NewMessagePublishHandler(bgService service.BackgroundJobServiceInterface, msgService service.MessageServiceInterface, redisMsgService service.RedisMessageServiceInterface, httpClient service.HttpClientInterface, lgr *zap.SugaredLogger) *MessagePublishHandler {
	return &MessagePublishHandler{
		BackgroundService:   bgService,
		MessageService:      msgService,
		RedisMessageService: redisMsgService,
		HttpClient:          httpClient,
		logger:              lgr,
	}
}
