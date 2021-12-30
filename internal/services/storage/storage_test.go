package storage

import (
	"context"
	"errors"
	"fmt"
	"github.com/BrunoDM2943/church-members-api/internal/constants/domain"
	mock_wrapper "github.com/BrunoDM2943/church-members-api/platform/aws/wrapper/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestS3Storage_SaveFile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var churchID = domain.NewID()
	const fileName = "members_report.pdf"
	var data = []byte("my-file")
	var key = fmt.Sprintf("%s/%s", churchID, fileName)

	s3Wrapper := mock_wrapper.NewMockS3APIWrapper(ctrl)
	storage := NewS3Storage(s3Wrapper)

	t.Run("Save file successfully", func(t *testing.T) {
		ctx := BuildContext(churchID)
		s3Wrapper.EXPECT().PutObject(gomock.Eq(ctx), gomock.Eq(key), gomock.Eq(data)).Return(nil)
		assert.Nil(t, storage.SaveFile(ctx, fileName, data))
	})
	t.Run("Save file returns error", func(t *testing.T) {
		ctx := BuildContext(churchID)
		s3Wrapper.EXPECT().PutObject(gomock.Eq(ctx), gomock.Eq(key), gomock.Eq(data)).Return(errors.New("error"))
		assert.NotNil(t, storage.SaveFile(ctx, fileName, data))
	})
}

func BuildContext(id string) context.Context {
	return context.WithValue(context.TODO(), "user", &domain.User{
		Church: &domain.Church{
			ID: id,
		},
	})
}
