package api

import (
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/mburaksoran/insider-case/internal/app/service"
	"go.uber.org/zap"
)

type messageHandler struct {
	MessageService service.MessageServiceInterface
	logger         *zap.SugaredLogger
}

type MessageHandlerInterface interface {
	ListSendMessagesHandler(c *fiber.Ctx) error
}

func NewMessageHandler(messageService service.MessageServiceInterface, lgr *zap.SugaredLogger) MessageHandlerInterface {
	return &messageHandler{
		MessageService: messageService,
		logger:         lgr,
	}
}

func (mh *messageHandler) ListSendMessagesHandler(c *fiber.Ctx) error {

	messages, err := mh.MessageService.GetSendMessage(c.Context())
	if err != nil {
		mh.logger.Error("[ListSendMessagesHandler] - Error while listing messages, Error: ", err)
		return errors.New(fmt.Sprintf("[ListSendMessagesHandler] -  Error while listing messages, Error: %v", err))
	}

	return c.JSON(messages)
}
