package wrapper

import (
	"bytes"
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/sirupsen/logrus"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

//go:generate mockgen -source=./s3.go -destination=./mock/s3_mock.go
type s3api interface {
	PutObject(ctx context.Context,
		params *s3.PutObjectInput,
		optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error)
}

type S3APIWrapper interface {
	PutObject(ctx context.Context, key string, data []byte) error
}

type s3wrapper struct {
	api    s3api
	bucket string
}

func NewS3APIWrapper(api s3api, bucket string) S3APIWrapper {
	return &s3wrapper{api: api, bucket: bucket}
}

func (wrapper *s3wrapper) PutObject(ctx context.Context, key string, data []byte) error {
	input := &s3.PutObjectInput{
		Bucket: aws.String(wrapper.bucket),
		Key:    aws.String(key),
		Body:   bytes.NewReader(data),
	}

	_, err := wrapper.api.PutObject(ctx, input)
	if err != nil {
		logrus.Errorf("Error uploading object to S3 bucket: %v", err)
		return err
	}
	return nil
}
