package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/shashank/order-service/internal/auth"
)

func TestAuthMiddleware(t *testing.T) {

	gin.SetMode(gin.TestMode)

	validToken, err := auth.GenerateToken("TEST_SHASHANK")

	if err != nil {
		t.Fatalf("failed to generate test token: %v", err)
	}

	tests := []struct {
		name            string
		authHeader      string
		expectedStatus  int
		expectedUser    string
		handlerExpected bool
	}{
		{
			name:            "Token is missing",
			authHeader:      "",
			expectedStatus:  http.StatusUnauthorized,
			handlerExpected: false,
		},
		{
			name:            "Token format is invalid, Bearer prefix is missing",
			authHeader:      "jfdksajk",
			expectedStatus:  http.StatusUnauthorized,
			handlerExpected: false,
		},
		{
			name:            "Invalid JWT",
			authHeader:      "Bearer jfdksajk",
			expectedStatus:  http.StatusUnauthorized,
			handlerExpected: false,
		},
		{
			name:            "Valid JWT",
			authHeader:      "Bearer " + validToken,
			expectedStatus:  http.StatusOK,
			expectedUser:    "TEST_SHASHANK",
			handlerExpected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			w := httptest.NewRecorder()

			handlerCalled := false
			var receivedUser string

			router := gin.New()

			router.Use(AuthMiddleware())

			router.POST("/orders", func(ctx *gin.Context) {

				handlerCalled = true

				user, exists := ctx.Get("userId")

				if exists {
					userString, ok := user.(string)

					if !ok {
						t.Errorf("Expected userId to be string, but got %T", user)
					} else {
						receivedUser = userString
					}
				}
			})

			req := httptest.NewRequest(
				http.MethodPost,
				"/orders",
				nil,
			)

			req.Header.Set("Authorization", tt.authHeader)

			router.ServeHTTP(w, req)

			// Check HTTP status
			if w.Code != tt.expectedStatus {
				t.Errorf(
					"Expected status %d but received %d",
					tt.expectedStatus,
					w.Code,
				)
			}

			// Check whether next handler was executed
			if handlerCalled != tt.handlerExpected {
				t.Errorf(
					"Expected handler called = %v but got %v",
					tt.handlerExpected,
					handlerCalled,
				)
			}

			// Check userId
			if tt.expectedUser != "" {
				if receivedUser != tt.expectedUser {
					t.Errorf(
						"Expected user %s but received %s",
						tt.expectedUser,
						receivedUser,
					)
				}
			}
		})
	}
}
