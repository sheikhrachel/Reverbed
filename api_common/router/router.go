package router

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	timeout "github.com/s-wijaya/gin-timeout"

	"github.com/sheikhrachel/reverbed/api_common/call"
	"github.com/sheikhrachel/reverbed/api_common/utils/errutil"
)

const (
	timeoutVal        = 60 * time.Second
	errRequestTimeout = "request timeout"
)

var (
	corsConfig = cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false,
	}
)

// InitCall is a function that will set up the call context for the application using
// the ENV and REGION environment variables
func InitCall() (cc call.Call) {
	// set the env
	appEnv := os.Getenv("ENV")
	if appEnv == "" {
		appEnv = "local"
		os.Setenv("ELASTICACHE_ENDPOINT", "localhost:11211")
		defer os.Unsetenv("ELASTICACHE_ENDPOINT")
	}
	// set the region
	appRegion := os.Getenv("REGION")
	if appRegion == "" {
		appRegion = "us-west-2"
	}
	// setting up the call context
	// - this is a custom package that will help us create distributed traces
	return call.New(appEnv, appRegion)
}

// CreateRouterWithMiddleware is a function that will create a new gin router
// and set up the general middleware that we need for the application
func CreateRouterWithMiddleware() (router *gin.Engine) {
	// setting the gin router to release mode
	// - this will disable the debug mode, which has overtly verbose console
	//   logging that we don't need in deployed instances
	gin.SetMode(gin.ReleaseMode)
	router = gin.Default()
	// setting up the cors middleware
	// - this is a custom middleware that will set the cors headers
	SetCorsOnRouter(router)
	// setting up the timeout middleware
	// - this is a custom middleware that will time out incoming requests
	//   if it takes longer than 40 seconds
	SetupTimeoutMiddleware(router)
	return router
}

// StartRouter starts the router on port 8080
func StartRouter(r *gin.Engine, cc call.Call) {
	cc.InfoF("router starting on 8080, env: %s", cc.Env)
	err := r.Run("0.0.0.0:8080")
	errutil.HandleError(cc, err)
}

// SetCorsOnRouter sets up the cors configuration on the router
func SetCorsOnRouter(r *gin.Engine) {
	r.Use(cors.New(corsConfig))
}

// SetupTimeoutMiddleware sets up the timeout middleware functionality on the router
func SetupTimeoutMiddleware(r *gin.Engine) {
	r.Use(timeout.TimeoutHandler(timeoutVal, http.StatusRequestTimeout, errRequestTimeout))
}
