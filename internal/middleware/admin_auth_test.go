package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestAdminBearerAuth(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name         string
		rootToken    string
		authHeader   string
		expectedCode int
	}{
		{"ValidToken", "secret123", "Bearer secret123", http.StatusOK},
		{"InvalidToken", "secret123", "Bearer wrong", http.StatusUnauthorized},
		{"MissingHeader", "secret123", "", http.StatusUnauthorized},
		{"WrongPrefix", "secret123", "Basic secret123", http.StatusUnauthorized},
		{"EmptyConfig", "", "Bearer anything", http.StatusUnauthorized},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()
			router.Use(AdminBearerAuth(tt.rootToken))
			router.GET("/test", func(c *gin.Context) {
				c.Status(http.StatusOK)
			})

			req, _ := http.NewRequest("GET", "/test", nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)
		})
	}
}
