package handlers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/alonsoF100/chat-api/internal/models"
	"github.com/alonsoF100/chat-api/internal/transport/http/dto"
	"github.com/alonsoF100/chat-api/internal/transport/http/handlers"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateChat(t *testing.T) {
	mc := minimock.NewController(t)
	mockService := handlers.NewServiceMock(mc)
	handler := handlers.New(mockService)

	tests := []struct {
		name           string
		requestBody    string
		setupMocks     func()
		expectedStatus int
		expectedError  string
	}{
		{
			name:        "successful creation",
			requestBody: `{"title": "Test Chat"}`,
			setupMocks: func() {
				expectedChat := &models.Chat{
					ID:    1,
					Title: "Test Chat",
				}
				mockService.CreateChatMock.Expect("Test Chat").Return(expectedChat, nil)
			},
			expectedStatus: http.StatusCreated,
			expectedError:  "",
		},

		{
			name:           "invalid json format",
			requestBody:    `{"title": "Test Chat",}`,
			setupMocks:     func() {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "invalid character",
		},

		{
			name:           "unknown fields",
			requestBody:    `{"title": "Test Chat", "unknown": "field"}`,
			setupMocks:     func() {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "unknown field",
		},

		{
			name:           "missing title field",
			requestBody:    `{}`,
			setupMocks:     func() {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "required",
		},

		{
			name:           "empty title",
			requestBody:    `{"title": ""}`,
			setupMocks:     func() {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "required",
		},

		{
			name:           "invalid title type",
			requestBody:    `{"title": 123}`,
			setupMocks:     func() {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "cannot unmarshal number",
		},

		{
			name:           "title too long",
			requestBody:    `{"title": "` + strings.Repeat("a", 256) + `"}`,
			setupMocks:     func() {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "max",
		},

		{
			name:        "service error",
			requestBody: `{"title": "Test Chat"}`,
			setupMocks: func() {
				mockService.CreateChatMock.Expect("Test Chat").Return(nil, errors.New("database connection failed"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  "database connection failed",
		},

		{
			name:        "title with spaces",
			requestBody: `{"title": "  Test Chat  "}`,
			setupMocks: func() {
				expectedChat := &models.Chat{
					ID:    2,
					Title: "Test Chat",
				}
				mockService.CreateChatMock.Expect("  Test Chat  ").Return(expectedChat, nil)
			},
			expectedStatus: http.StatusCreated,
			expectedError:  "",
		},

		{
			name:           "empty request body",
			requestBody:    "",
			setupMocks:     func() {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "EOF",
		},

		{
			name:        "special characters in title",
			requestBody: `{"title": "Chat #1 @test"}`,
			setupMocks: func() {
				expectedChat := &models.Chat{
					ID:    3,
					Title: "Chat #1 @test",
				}
				mockService.CreateChatMock.Expect("Chat #1 @test").Return(expectedChat, nil)
			},
			expectedStatus: http.StatusCreated,
			expectedError:  "",
		},

		{
			name:        "minimum length title",
			requestBody: `{"title": "a"}`,
			setupMocks: func() {
				expectedChat := &models.Chat{
					ID:    4,
					Title: "a",
				}
				mockService.CreateChatMock.Expect("a").Return(expectedChat, nil)
			},
			expectedStatus: http.StatusCreated,
			expectedError:  "",
		},

		{
			name:        "maximum length title",
			requestBody: `{"title": "` + strings.Repeat("a", 200) + `"}`,
			setupMocks: func() {
				expectedChat := &models.Chat{
					ID:    5,
					Title: strings.Repeat("a", 200),
				}
				mockService.CreateChatMock.Expect(strings.Repeat("a", 200)).Return(expectedChat, nil)
			},
			expectedStatus: http.StatusCreated,
			expectedError:  "",
		},

		{
			name:        "multiline title",
			requestBody: `{"title": "Line 1\nLine 2"}`,
			setupMocks: func() {
				expectedChat := &models.Chat{
					ID:    6,
					Title: "Line 1\nLine 2",
				}
				mockService.CreateChatMock.Expect("Line 1\nLine 2").Return(expectedChat, nil)
			},
			expectedStatus: http.StatusCreated,
			expectedError:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			req := httptest.NewRequest(http.MethodPost, "/chats/", bytes.NewBufferString(tt.requestBody))
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()
			handler.CreateChat(rr, req)

			require.Equal(t, tt.expectedStatus, rr.Code)

			if tt.expectedError != "" {
				var errorResp dto.ErrorResponse

				err := json.Unmarshal(rr.Body.Bytes(), &errorResp)
				require.NoError(t, err)
				
				assert.Contains(t, errorResp.Error, tt.expectedError)
				assert.NotEmpty(t, errorResp.Timestamp)
			} else if tt.expectedStatus == http.StatusCreated {
				
				var chatResp dto.ChatResponse
				err := json.Unmarshal(rr.Body.Bytes(), &chatResp)
				require.NoError(t, err)
				assert.NotZero(t, chatResp.ID)
				assert.NotEmpty(t, chatResp.Title)
			}
		})
	}
}
