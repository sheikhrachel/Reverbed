package secrets_manager

import (
	"errors"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/secretsmanager"

	"github.com/sheikhrachel/reverbed/api_common/call"
	"github.com/sheikhrachel/reverbed/api_common/utils/errutil"
)

// GetSecretValue gets a corresponding value string for a key in secrets manager
func (s *Service) GetSecretValue(cc call.Call, key string) (string, error) {
	if val := SecretsCache.Get(key); val != nil {
		return val.Value(), nil
	}
	result, err := s.Svc.GetSecretValue(&secretsmanager.GetSecretValueInput{SecretId: aws.String(key)})
	if errutil.HandleError(cc, err) {
		handleSecretsManagerErr(cc, err)
		return "", err
	}
	if result.SecretString == nil {
		return "", nil
	}
	SecretsCache.Set(key, *(result.SecretString), 30*time.Minute)
	return *(result.SecretString), nil
}

// handleSecretsManagerErr handles errors from Secrets Manager
func handleSecretsManagerErr(cc call.Call, err error) {
	if err == nil {
		return
	}
	var (
		errMsg, errStr string
		aErr           awserr.Error
	)
	if errors.As(err, &aErr) {
		errStr, errMsg = aErr.Error(), aErr.Code()
	}
	cc.InfoF("msg: %#v, err: %#v", errMsg, errStr)
}
