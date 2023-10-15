package dynamodb

import (
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"

	"github.com/sheikhrachel/reverbed/api_common/call"
)

type Service struct {
	Client
	Svc dynamodbiface.DynamoDBAPI
}

func InitDynamo(cc call.Call) *Service {
	awsConfig := &aws.Config{Region: aws.String(cc.Region)}
	if cc.Env == "local" {
		accessKey, secretKey := os.Getenv("AWS_ACCESS_KEY_ID"), os.Getenv("AWS_SECRET_ACCESS_KEY")
		if accessKey == "" || secretKey == "" {
			accessKey, secretKey = "test", "test"
		}
		cc.TraceF("Using local DynamoDB instance with accessKey: %s, secretKey: %s", accessKey, secretKey)
		awsConfig.Endpoint = aws.String("http://local-dynamodb:8000")
		awsConfig.Credentials = credentials.NewStaticCredentials(accessKey, secretKey, "")
	}
	svc := dynamodb.New(session.Must(session.NewSessionWithOptions(session.Options{SharedConfigState: session.SharedConfigEnable})), awsConfig)
	cc.InfoF("Successfully initialized DynamoDB client")
	return &Service{Svc: svc}
}
