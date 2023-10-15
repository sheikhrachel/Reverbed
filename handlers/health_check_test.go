package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/smarty/assertions/should"

	"github.com/sheikhrachel/reverbed/api_common/utils/testutil"
)

type MockHandler struct {
	Handler
}

func TestHealthCheck(t *testing.T) {
	handler := &MockHandler{}
	gin.SetMode(gin.TestMode)
	testCases := []struct {
		name, path, expectedBody string
		expectedStatus           int
	}{
		{"health check", "/health", "\"health\"", http.StatusOK},
		{"root check", "/", "\"health\"", http.StatusOK},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			router := gin.Default()
			router.GET(tc.path, handler.HealthCheck)
			req, err := http.NewRequest(http.MethodGet, tc.path, nil)
			testutil.So(t, err, should.BeNil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			testutil.So(t, w.Code, should.Equal, tc.expectedStatus)
			testutil.So(t, w.Body.String(), should.Equal, tc.expectedBody)
		})
	}
}
