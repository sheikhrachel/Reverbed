package secrets_manager

import (
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"

	"github.com/sheikhrachel/reverbed/api_common/call"
	"github.com/sheikhrachel/reverbed/api_common/utils/errutil"
)

type Service struct {
	Svc *secretsmanager.SecretsManager
}

// SecretsCache is a cache of secrets from Secrets Manager and is used to avoid
// making multiple calls to Secrets Manager for the same secret.
var SecretsCache = ttlcache.New[string, string](ttlcache.WithTTL[string, string](30 * time.Minute))

func InitSM(cc call.Call) (service *Service) {
	// Create a Secrets Manager client
	sess, err := session.NewSession(&aws.Config{Region: aws.String(cc.Region)})
	if errutil.HandleError(cc, err) {
		// Handle session creation error
		return nil
	}
	return &Service{Svc: secretsmanager.New(sess)}
}
