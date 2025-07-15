package yago

import (
	"net/http"
	"time"
)

// Client is the main entry point for interacting with the YandexGPT API.
// It manages authentication and request sending.

const (
	baseURL = "https://llm.api.cloud.yandex.net/foundationModels/v1/completion"
)

type Client struct {
	httpClient *http.Client
	baseURL    string
	apiKey     string
	folderID   string
}

// Option is a function type used to configure a Client.
type Option func(*Client)

// WithCustomURL provides an Option to set a custom base URL for the API.
func WithCustomURL(url string) Option {
	return func(c *Client) {
		c.baseURL = url
	}
}

// WithCustomClient provides an Option to use a custom http.Client.
func WithCustomClient(client *http.Client) Option {
	return func(c *Client) {
		c.httpClient = client
	}
}

// WithCustomTimeout provides an Option to set a custom timeout for the http.Client.
func WithCustomTimeout(timeout time.Duration) Option {
	return func(c *Client) {
		if c.httpClient == nil {
			c.httpClient = &http.Client{}
		}
		c.httpClient.Timeout = timeout
	}
}

// NewClient creates, configures, and returns a new Client instance.
func NewClient(apiKey string, folderID string, opts ...Option) *Client {
	client := &Client{
		apiKey: apiKey,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		baseURL:  baseURL,
		folderID: folderID,
	}

	for _, opt := range opts {
		opt(client)
	}

	return client
}

func (c *Client) Close() {
	c.httpClient.CloseIdleConnections()
}
