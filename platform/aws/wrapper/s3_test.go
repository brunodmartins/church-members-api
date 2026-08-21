package wrapper

import (
	"context"
	"errors"
	"testing"
	"time"

	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	mock_wrapper "github.com/brunodmartins/church-members-api/platform/aws/wrapper/mock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestS3wrapper_PutObject(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	const fileName = "members_report.pdf"
	const bucket = "my-bucket"
	var data = []byte("my-file")

	api := mock_wrapper.NewMocks3API(ctrl)
	wrapper := NewS3APIWrapper(api, bucket, nil)
	ctx := context.TODO()

	t.Run("Success", func(t *testing.T) {
		api.EXPECT().PutObject(gomock.Eq(ctx), gomock.Any()).DoAndReturn(func(ctx context.Context,
			params *s3.PutObjectInput,
			optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error) {
			assert.Equal(t, bucket, *params.Bucket)
			assert.Equal(t, fileName, *params.Key)
			assert.NotNil(t, params.Body)
			return nil, nil
		})
		assert.Nil(t, wrapper.PutObject(ctx, fileName, data))
	})
	t.Run("Fail", func(t *testing.T) {
		api.EXPECT().PutObject(gomock.Eq(ctx), gomock.Any()).DoAndReturn(func(ctx context.Context,
			params *s3.PutObjectInput,
			optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error) {
			assert.Equal(t, bucket, *params.Bucket)
			assert.Equal(t, fileName, *params.Key)
			assert.NotNil(t, params.Body)
			return nil, errors.New("error")
		})
		assert.NotNil(t, wrapper.PutObject(ctx, fileName, data))
	})
}

func TestS3wrapper_PresignGetObject(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	const fileName = "members_report.pdf"
	const bucket = "my-bucket"
	const signedURL = "signed-url"

	api := mock_wrapper.NewMocks3SignedAPI(ctrl)
	wrapper := NewS3APIWrapper(nil, bucket, api)
	ctx := context.TODO()

	t.Run("Success", func(t *testing.T) {
		api.EXPECT().PresignGetObject(gomock.Eq(ctx), gomock.Any()).DoAndReturn(func(ctx context.Context,
			params *s3.GetObjectInput,
			optFns ...func(*s3.Options)) (*v4.PresignedHTTPRequest, error) {
			assert.Equal(t, bucket, *params.Bucket)
			assert.Equal(t, fileName, *params.Key)
			return &v4.PresignedHTTPRequest{
				URL: signedURL,
			}, nil
		})
		result, err := wrapper.PresignGetObject(ctx, fileName)
		assert.Nil(t, err)
		assert.Equal(t, signedURL, result)
	})
	t.Run("Fail", func(t *testing.T) {
		api.EXPECT().PresignGetObject(gomock.Eq(ctx), gomock.Any()).DoAndReturn(func(ctx context.Context,
			params *s3.GetObjectInput,
			optFns ...func(*s3.Options)) (*v4.PresignedHTTPRequest, error) {
			assert.Equal(t, bucket, *params.Bucket)
			assert.Equal(t, fileName, *params.Key)
			return nil, errors.New("error")
		})
		_, err := wrapper.PresignGetObject(ctx, fileName)
		assert.NotNil(t, err)
	})
}

func TestS3wrapper_HeadObject(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	const fileName = "members_report.pdf"
	const bucket = "my-bucket"

	api := mock_wrapper.NewMocks3API(ctrl)
	wrapper := NewS3APIWrapper(api, bucket, nil)
	ctx := context.TODO()

	t.Run("Success", func(t *testing.T) {
		api.EXPECT().HeadObject(gomock.Eq(ctx), gomock.Any()).DoAndReturn(func(ctx context.Context,
			params *s3.HeadObjectInput,
			optFns ...func(*s3.Options)) (*s3.HeadObjectOutput, error) {
			assert.Equal(t, bucket, *params.Bucket)
			assert.Equal(t, fileName, *params.Key)
			return &s3.HeadObjectOutput{
				LastModified:  func() *time.Time { value := time.Now(); return &value }(),
				ContentType:   func() *string { value := "application/pdf"; return &value }(),
				ContentLength: func() *int64 { value := int64(42); return &value }(),
			}, nil
		})
		result, err := wrapper.HeadObject(ctx, fileName)
		assert.Nil(t, err)
		assert.Equal(t, "application/pdf", result["ContentType"])
		assert.Equal(t, int64(42), result["Size"])
	})
	t.Run("Fail", func(t *testing.T) {
		api.EXPECT().HeadObject(gomock.Eq(ctx), gomock.Any()).DoAndReturn(func(ctx context.Context,
			params *s3.HeadObjectInput,
			optFns ...func(*s3.Options)) (*s3.HeadObjectOutput, error) {
			assert.Equal(t, bucket, *params.Bucket)
			assert.Equal(t, fileName, *params.Key)
			return nil, errors.New("error")
		})
		_, err := wrapper.HeadObject(ctx, fileName)
		assert.NotNil(t, err)
	})
}
