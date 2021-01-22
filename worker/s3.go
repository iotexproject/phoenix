package worker

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// import (
// 	"context"
// 	"time"

// 	"github.com/nsqio/go-diskqueue"
// )

// type worker struct {
// 	cb       func()
// 	duration time.Duration
// 	ch       chan interface{}
// 	queue    diskqueue.Interface
// }

func S3(job Job) {

}

func createS3SessionByAcessKey(accessKey string, accessToken string, region string) (*session.Session, error) {
	return session.NewSession(&aws.Config{
		Region: aws.String(region),
		Credentials: credentials.NewStaticCredentials(
			accessKey,
			accessToken,
			""),
	})
}

func NewAmazonS3WithCredentials(region string, endpoint string, credentials *credentials.Credentials) *s3.S3 {
	return s3.New(session.New(), &aws.Config{
		Credentials:      credentials,
		Region:           aws.String(region),
		Endpoint:         aws.String(endpoint),
		DisableSSL:       aws.Bool(strings.HasPrefix(endpoint, "http://")),
		S3ForcePathStyle: aws.Bool(endpoint != ""),
	})

func createAWSConfigByARN(arn string, externalID string, region string, sess *session.Session) aws.Config {

	conf := aws.Config{Region: aws.String(region)}
	if arn != "" {
		var creds *credentials.Credentials
		if externalID != "" {
			creds = stscreds.NewCredentials(sess, arn, func(p *stscreds.AssumeRoleProvider) {
				p.ExternalID = &externalID
			})
		} else {
			creds = stscreds.NewCredentials(sess, arn, func(p *stscreds.AssumeRoleProvider) {})
		}
		conf.Credentials = creds
	}
	return conf
}
