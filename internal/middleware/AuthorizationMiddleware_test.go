package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestAuthorizationMiddleWare(t *testing.T) {

	tests := []struct {
		name            string
		requiredRole    string
		status          int
		roleToGenerate  any
		isHandlerCalled bool
		expectedError   string
		isErrorExpected bool
	}{
		{
			name:            "Roles is missing",
			requiredRole:    "TEST_ROLE",
			status:          http.StatusForbidden,
			isHandlerCalled: false,
			expectedError:   "User not authorized",
			isErrorExpected: true,
		},
		{
			name:            "Role type mismatch",
			requiredRole:    "TEST_ROLE",
			status:          http.StatusForbidden,
			roleToGenerate:  []int{1, 2, 3},
			isHandlerCalled: false,
			expectedError:   "Invalid user roles",
			isErrorExpected: true,
		},
		{
			name:            "Required role is missing",
			requiredRole:    "TEST_ROLE",
			status:          http.StatusForbidden,
			roleToGenerate:  []string{"TEST_USR", "TST_USER"},
			isHandlerCalled: false,
			expectedError:   "Required user role not found",
			isErrorExpected: true,
		},
		{
			name:            "Required role is Present",
			requiredRole:    "TEST_ROLE",
			status:          http.StatusCreated,
			roleToGenerate:  []string{"TEST_USR", "TEST_ROLE"},
			isHandlerCalled: true,
			isErrorExpected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			handlerCalled := false

			router := gin.New()

			router.Use(func(ctx *gin.Context) {
				ctx.Set("roles", tt.roleToGenerate)
				ctx.Next()
			})
			router.Use(RequiredRole([]string{tt.requiredRole}))

			router.POST("/orders", func(ctx *gin.Context) {
				handlerCalled = true
				ctx.JSON(http.StatusCreated, nil)
			})

			req := httptest.NewRequest(http.MethodPost, "/orders", nil)

			router.ServeHTTP(w, req)

			if w.Code != tt.status {
				t.Errorf("Exepected status was %v but received %v", tt.status, w.Code)
			}

			if tt.isHandlerCalled != handlerCalled {
				t.Fatalf("Has hander to be called ?? %v, In real is it called?? %v ", tt.isHandlerCalled, handlerCalled)
			}

		})
	}
}
