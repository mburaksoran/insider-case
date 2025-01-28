package service

import (
	"context"
)

type HttpClientInterface interface {
	PostWithAPIKey(ctx context.Context, data interface{}) ([]byte, error)
}
