package handlers

import (
	"github.com/sheikhrachel/reverbed/reverbed_service"

	aws "github.com/sheikhrachel/reverbed/api_common/aws_utils"
	"github.com/sheikhrachel/reverbed/api_common/call"
)

// Handler is a struct that holds the common dependencies for all endpoint handlers
type Handler struct {
	// appEnv is the environment the app is running in
	appEnv string
	// appRegion is the region the app is running in
	appRegion string
	// aws is the aws client
	aws *aws.AWSClient
	// reverbedService is the service interface responsible for business logic
	reverbedService *reverbed_service.ReverbedService
}

// New returns a new Handler pointer
func New(cc call.Call) *Handler {
	awsClient := aws.Init(cc)
	return &Handler{
		appEnv:          cc.Env,
		appRegion:       cc.Region,
		aws:             awsClient,
		reverbedService: reverbed_service.New(cc, awsClient),
	}
}
