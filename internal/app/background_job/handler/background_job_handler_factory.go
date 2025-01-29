package handler

import (
	"context"
	"fmt"
	"github.com/mburaksoran/insider-case/internal/app/service"
	"github.com/mburaksoran/insider-case/internal/app/types"
	"go.uber.org/zap"
)

type JobHandler interface {
	Handle(ctx context.Context, values interface{}) error
}

type HandlerFactory struct {
	BackgroundService   service.BackgroundJobServiceInterface
	MessageService      service.MessageServiceInterface
	RedisMessageService service.RedisMessageServiceInterface
	HttpClient          service.HttpClientInterface
	logger              *zap.SugaredLogger
}

func NewHandlerFactory(bgService service.BackgroundJobServiceInterface, msgService service.MessageServiceInterface, redisMsgService service.RedisMessageServiceInterface, httpClient service.HttpClientInterface, lgr *zap.SugaredLogger) *HandlerFactory {
	return &HandlerFactory{
		BackgroundService:   bgService,
		MessageService:      msgService,
		RedisMessageService: redisMsgService,
		HttpClient:          httpClient,
		logger:              lgr,
	}
}

func (f *HandlerFactory) GetHandler(handlerName string) (JobHandler, error) {
	switch handlerName {
	case types.MessagePublishHandler:
		return NewMessagePublishHandler(f.BackgroundService, f.MessageService, f.RedisMessageService, f.HttpClient, f.logger), nil
	default:
		return nil, fmt.Errorf("handler not found: %s", handlerName)
	}
}
