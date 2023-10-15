package sqs

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"

	"github.com/sheikhrachel/reverbed/api_common/call"
)

type Service struct {
	Svc *sqs.SQS
}

func InitSQS(cc call.Call) *Service {
	return &Service{
		Svc: sqs.New(session.Must(session.NewSession(&aws.Config{Region: aws.String(cc.Region)}))),
	}
}
