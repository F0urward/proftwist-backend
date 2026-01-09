package baseclient

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/F0urward/proftwist-backend/pkg/ctxutil"
)

type BaseClient struct {
	httpClient *http.Client
}

func NewBaseClient() *BaseClient {
	return &BaseClient{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func NewInsecureBaseClient() *BaseClient {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	return &BaseClient{
		httpClient: &http.Client{
			Timeout:   60 * time.Second,
			Transport: transport,
		},
	}
}

func (c *BaseClient) DoRequest(ctx context.Context, req *http.Request, result interface{}) error {
	const op = "BaseClient.DoRequest"
	logger := ctxutil.GetLogger(ctx).WithField("op", op)

	req = req.WithContext(ctx)

	logger.WithFields(map[string]interface{}{
		"method": req.Method,
		"url":    req.URL.String(),
	}).Info("making HTTP request")

	httpResp, err := c.httpClient.Do(req)
	if err != nil {
		logger.WithError(err).Error("HTTP request failed")
		return fmt.Errorf("HTTP request failed: %w", err)
	}
	defer func() {
		if err := httpResp.Body.Close(); err != nil {
			logger.WithError(err).Warn("failed to close response body")
		}
	}()

	if httpResp.StatusCode < 200 || httpResp.StatusCode > 299 {
		bodyBytes, _ := io.ReadAll(httpResp.Body)
		logger.WithFields(map[string]interface{}{
			"status_code": httpResp.StatusCode,
			"body":        string(bodyBytes),
		}).Error("server returned error status")
		return fmt.Errorf("server returned status %d: %s", httpResp.StatusCode, string(bodyBytes))
	}

	if err := json.NewDecoder(httpResp.Body).Decode(result); err != nil {
		logger.WithError(err).Error("failed to decode JSON response")
		return fmt.Errorf("failed to decode JSON response: %w", err)
	}

	logger.Info("successfully processed HTTP response")
	return nil
}
