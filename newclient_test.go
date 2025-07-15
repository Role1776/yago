package yago

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewClient(t *testing.T) {

	customClientInstance := &http.Client{
		Timeout:   100 * time.Second,
		Transport: &http.Transport{},
	}

	testCases := []struct {
		name                   string
		apiKey                 string
		folderID               string
		opts                   []Option
		expectedBaseURL        string
		expectedTimeout        time.Duration
		expectedClientInstance *http.Client
	}{
		{
			name:                   "Client creation with default settings",
			apiKey:                 "default-api-key",
			folderID:               "default-folder-id",
			opts:                   nil,
			expectedBaseURL:        baseURL,
			expectedTimeout:        30 * time.Second,
			expectedClientInstance: nil,
		},
		{
			name:                   "With option WithCustomURL",
			apiKey:                 "custom-url-key",
			folderID:               "custom-url-folder",
			opts:                   []Option{WithCustomURL("http://my.custom.url/api")},
			expectedBaseURL:        "http://my.custom.url/api",
			expectedTimeout:        30 * time.Second,
			expectedClientInstance: nil,
		},
		{
			name:                   "With option WithCustomTimeout",
			apiKey:                 "custom-timeout-key",
			folderID:               "custom-timeout-folder",
			opts:                   []Option{WithCustomTimeout(15 * time.Second)},
			expectedBaseURL:        baseURL,
			expectedTimeout:        15 * time.Second,
			expectedClientInstance: nil,
		},
		{
			name:                   "With option WithCustomClient",
			apiKey:                 "custom-client-key",
			folderID:               "custom-client-folder",
			opts:                   []Option{WithCustomClient(customClientInstance)},
			expectedBaseURL:        baseURL,
			expectedTimeout:        customClientInstance.Timeout,
			expectedClientInstance: customClientInstance,
		},
		{
			name:     "With multiple options (URL and Timeout)",
			apiKey:   "multi-option-key",
			folderID: "multi-option-folder",
			opts: []Option{
				WithCustomURL("http://another.url/v2"),
				WithCustomTimeout(5 * time.Second),
			},
			expectedBaseURL:        "http://another.url/v2",
			expectedTimeout:        5 * time.Second,
			expectedClientInstance: nil,
		},
		{
			name:     "Option WithCustomTimeout overrides timeout in WithCustomClient",
			apiKey:   "override-key",
			folderID: "override-folder",
			opts: []Option{
				WithCustomClient(customClientInstance),
				WithCustomTimeout(25 * time.Second),
			},
			expectedBaseURL:        baseURL,
			expectedTimeout:        25 * time.Second,
			expectedClientInstance: customClientInstance,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			client := NewClient(tc.apiKey, tc.folderID, tc.opts...)

			require.NotNil(t, client, "NewClient should not return nil")
			require.NotNil(t, client.httpClient, "httpClient should not be nil after creation")

			assert.Equal(t, tc.apiKey, client.apiKey, "API key should be set correctly")
			assert.Equal(t, tc.folderID, client.folderID, "FolderID should be set correctly")
			assert.Equal(t, tc.expectedBaseURL, client.baseURL, "BaseURL should be correct")
			assert.Equal(t, tc.expectedTimeout, client.httpClient.Timeout, "Timeout should be correct")

			if tc.expectedClientInstance != nil {
				assert.Same(t, tc.expectedClientInstance, client.httpClient, "The same http.Client instance should be used")
			} else {
				assert.NotSame(t, customClientInstance, client.httpClient, "Custom client should not be used if not provided")
			}
		})
	}
}
