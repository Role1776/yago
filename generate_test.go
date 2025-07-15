package yago

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerativeModel_Generate(t *testing.T) {
	testCases := []struct {
		name              string
		messages          []Message
		mockAPIHandler    http.HandlerFunc
		expectedResponse  *Response
		expectError       bool
		expectedErrorText string
	}{
		{
			name:     "Success",
			messages: []Message{{Role: "user", Text: "Привет, мир!"}},
			mockAPIHandler: func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "POST", r.Method)
				assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
				assert.Equal(t, "Api-Key super-secret-key", r.Header.Get("Authorization"))
				var payload requestPayload
				err := json.NewDecoder(r.Body).Decode(&payload)
				require.NoError(t, err)
				assert.Equal(t, "gpt://test-folder/yandexgpt-lite", payload.ModelURI)
				assert.Len(t, payload.Messages, 1)
				assert.Equal(t, "Привет, мир!", payload.Messages[0].Text)

				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{
					"result": {
						"alternatives": [
							{
								"message": {
									"role": "assistant",
									"text": "И тебе привет!"
								},
								"status": "ALTERNATIVE_STATUS_FINAL"
							}
						]
					}
				}`))
			},
			expectedResponse: &Response{
				Result: Result{
					Alternatives: []Alternative{
						{
							Message: Message{Role: "assistant", Text: "И тебе привет!"},
							Status:  "ALTERNATIVE_STATUS_FINAL",
						},
					},
				},
			},
			expectError: false,
		},
		{
			name:              "Empty message list",
			messages:          []Message{},
			mockAPIHandler:    nil,
			expectError:       true,
			expectedErrorText: "empty message list",
		},
		{
			name:     "Server error 500",
			messages: []Message{{Role: "user", Text: "Сломайся"}},
			mockAPIHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("internal server error"))
			},
			expectError:       true,
			expectedErrorText: "unexpected status code 500: internal server error",
		},
		{
			name:     "Server error 401 Unauthorized",
			messages: []Message{{Role: "user", Text: "Кто я?"}},
			mockAPIHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("api key is invalid"))
			},
			expectError:       true,
			expectedErrorText: "unexpected status code 401: api key is invalid",
		},
		{
			name:     "Invalid JSON response",
			messages: []Message{{Role: "user", Text: "Дай мне кривой JSON"}},
			mockAPIHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"result": "this is not valid`))
			},
			expectError:       true,
			expectedErrorText: "failed to decode response",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(tc.mockAPIHandler)
			defer server.Close()
			client := &Client{
				httpClient: server.Client(),
				baseURL:    server.URL,
				apiKey:     "super-secret-key",
				folderID:   "test-folder",
			}
			gm := &GenerativeModel{
				c:                 client,
				modelURI:          fmt.Sprintf("gpt://%s/yandexgpt-lite", client.folderID),
				CompletionOptions: CompletionOptions{},
			}

			response, err := gm.Generate(context.Background(), tc.messages)

			if tc.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedErrorText)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expectedResponse, response)
			}
		})
	}
}
