package s3

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"

	"github.com/sheikhrachel/reverbed/api_common/call"
	"github.com/sheikhrachel/reverbed/api_common/utils/errutil"
)

type Service struct {
	Downloader *s3manager.Downloader
	Svc        *s3.S3
	Uploader   *s3manager.Uploader
}

func InitS3(cc call.Call) *Service {
	sess, err := session.NewSession(&aws.Config{Region: aws.String(cc.Region)})
	if errutil.HandleError(cc, err) {
		return nil
	}
	return &Service{
		Downloader: s3manager.NewDownloader(sess),
		Svc:        s3.New(sess),
		Uploader:   s3manager.NewUploader(sess),
	}
}
