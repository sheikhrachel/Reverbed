package awsutils

import (
	"github.com/sheikhrachel/reverbed/api_common/call"
)

type AWSClient struct {
}

func Init(cc call.Call) (awsClient *AWSClient) {
	awsClient = &AWSClient{}
	return awsClient
}
