package chatapigateway

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

type Aws struct {
	Session *session.Session
}

func NewAwsSession(region, profile string) (*session.Session, error) {
	options := session.Options{
		Profile: profile, Config: aws.Config{
			Region: aws.String(region),
		},
	}
	return session.NewSessionWithOptions(options)
}
