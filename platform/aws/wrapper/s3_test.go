package wrapper

import (
	"context"
	"errors"
	mock_wrapper "github.com/BrunoDM2943/church-members-api/platform/aws/wrapper/mock"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestS3wrapper_PutObject(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	const fileName = "members_report.pdf"
	const bucket = "my-bucket"
	var data = []byte("my-file")

	api := mock_wrapper.NewMocks3api(ctrl)
	wrapper := NewS3APIWrapper(api, bucket)
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
