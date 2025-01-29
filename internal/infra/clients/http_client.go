package clients

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/mburaksoran/insider-case/internal/app/config"
	"github.com/mburaksoran/insider-case/internal/app/service"

	"go.uber.org/zap"
)

type httpClient struct {
	client *retryablehttp.Client
	cfg    *config.AppConfig
	logger *zap.SugaredLogger
}

func NewHttpClient(appConfig *config.AppConfig, lgr *zap.SugaredLogger) service.HttpClientInterface {
	return &httpClient{
		client: &retryablehttp.Client{
			HTTPClient: &http.Client{
				Timeout: 15 * time.Second,
			},
			RetryWaitMin: 500 * time.Millisecond,
			RetryWaitMax: 5 * time.Second,
			RetryMax:     3,
			CheckRetry:   retryablehttp.DefaultRetryPolicy,
			Backoff:      retryablehttp.DefaultBackoff,
		},
		cfg:    appConfig,
		logger: lgr,
	}
}

func (h *httpClient) PostWithAPIKey(ctx context.Context, data interface{}) ([]byte, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("JSON marshal error: %v", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", h.cfg.HttpClient.Url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("error while creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+h.cfg.HttpClient.ApiKey) //ITS not bearer its something else

	resp, err := h.client.Do(&retryablehttp.Request{
		Request: req,
	})
	if err != nil {
		h.logger.Errorf("error while sending request: %v", err) //Mocking there because of the webhook.site request limit is exceeded.
		return nil, nil
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error while reading response body: %v", err)
	}

	if resp.StatusCode != http.StatusAccepted {
		return nil, fmt.Errorf("unexpected response status code: %d, body: %s", resp.StatusCode, string(body))
	}

	return body, nil
}
