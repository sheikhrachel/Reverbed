package handlers

import (
	"expvar"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/sheikhrachel/reverbed/api_common/call"
)

func SetRoutes(router *gin.Engine, cc call.Call) (handler *Handler) {
	handler = New(cc)
	registerEndpoints(router, handler)
	return handler
}

// registerEndpoints acts as a registry for each route on the server
// {http method, path string, handler func}
func registerEndpoints(r *gin.Engine, h *Handler) {
	getEndpoints := []struct {
		method      string
		path        string
		handlerFunc gin.HandlerFunc
	}{
		// core
		{http.MethodGet, PathRoot, h.HealthCheck},
		{http.MethodGet, PathHealth, h.HealthCheck},
		// editing
		{http.MethodPost, PathEditTrack, h.EditTrack},
	}
	for _, endpoint := range getEndpoints {
		registerEndpoint(r, endpoint.method, endpoint.path, endpoint.handlerFunc)
	}
	// datadog expvars endpoint
	r.GET("/debug/vars", func(c *gin.Context) {
		expvar.Handler().ServeHTTP(c.Writer, c.Request)
	})
}

// registerEndpoint creates a new rate limiter for each path on the router, and
// assigns the route to the router with the correct HTTP method and handler func
func registerEndpoint(r *gin.Engine, method, path string, handlerFunc gin.HandlerFunc) {
	switch method {
	case http.MethodGet:
		r.GET(path, handlerFunc)
	case http.MethodPost:
		r.POST(path, handlerFunc)
	case http.MethodPut:
		r.PUT(path, handlerFunc)
	case http.MethodDelete:
		r.DELETE(path, handlerFunc)
	}
}
