package wrapper

import (
	"bytes"
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/sirupsen/logrus"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

//go:generate mockgen -source=./s3.go -destination=./mock/s3_mock.go
type s3API interface {
	PutObject(ctx context.Context,
		params *s3.PutObjectInput,
		optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error)
}

type s3SignedAPI interface {
	PresignGetObject(
		ctx context.Context,
		params *s3.GetObjectInput,
		optFns ...func(*s3.PresignOptions)) (*v4.PresignedHTTPRequest, error)
}

type S3APIWrapper interface {
	PutObject(ctx context.Context, key string, data []byte) error
	PresignGetObject(ctx context.Context, key string) (string, error)
}

type s3wrapper struct {
	api       s3API
	signedAPI s3SignedAPI
	bucket    string
}

func NewS3APIWrapper(api s3API, bucket string, s3SignedAPI s3SignedAPI) S3APIWrapper {
	return &s3wrapper{
		api:       api,
		bucket:    bucket,
		signedAPI: s3SignedAPI,
	}
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

func (wrapper *s3wrapper) PresignGetObject(ctx context.Context, key string) (string, error) {
	input := &s3.GetObjectInput{
		Bucket: aws.String(wrapper.bucket),
		Key:    aws.String(key),
	}

	request, err := wrapper.signedAPI.PresignGetObject(ctx, input)
	if err != nil {
		logrus.Errorf("Error generating signed url: %v", err)
		return "", err
	}
	return request.URL, nil
}
