package connection

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/koer/koer-module/pkg/config"
)

type RESTClient struct {
	baseURL string
	client  *http.Client
	headers map[string]string
}

func NewRESTClient(cfg config.RESTClientConfig) *RESTClient {
	timeout := cfg.Timeout
	if timeout == 0 {
		timeout = 30 * time.Second
	}
	return &RESTClient{
		baseURL: cfg.BaseURL,
		client:  &http.Client{Timeout: timeout},
		headers: cfg.Headers,
	}
}

func (c *RESTClient) do(ctx context.Context, method, path string, body any, result any) error {
	var bodyReader io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("marshaling request body: %w", err)
		}
		bodyReader = bytes.NewReader(b)
	}

	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, bodyReader)
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	for k, v := range c.headers {
		req.Header.Set(k, v)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	if result != nil {
		if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
			return fmt.Errorf("decoding response: %w", err)
		}
	}
	return nil
}

func (c *RESTClient) Get(ctx context.Context, path string, result any) error {
	return c.do(ctx, http.MethodGet, path, nil, result)
}

func (c *RESTClient) Post(ctx context.Context, path string, body, result any) error {
	return c.do(ctx, http.MethodPost, path, body, result)
}

func (c *RESTClient) Put(ctx context.Context, path string, body, result any) error {
	return c.do(ctx, http.MethodPut, path, body, result)
}

func (c *RESTClient) Delete(ctx context.Context, path string, result any) error {
	return c.do(ctx, http.MethodDelete, path, nil, result)
}
