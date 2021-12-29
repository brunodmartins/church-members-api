package storage

import (
	"context"
	"fmt"
	"github.com/BrunoDM2943/church-members-api/internal/constants/domain"
	"github.com/BrunoDM2943/church-members-api/platform/aws/wrapper"
)

//go:generate mockgen -source=./storage.go -destination=./mock/storage_mock.go
type Service interface {
	SaveFile(ctx context.Context, name string, data []byte) error
}

type s3Storage struct {
	s3Wrapper wrapper.S3APIWrapper
}

func NewS3Storage(s3API wrapper.S3APIWrapper) Service {
	return &s3Storage{s3Wrapper: s3API}
}

func (storage s3Storage) SaveFile(ctx context.Context, name string, data []byte) error {
	return storage.s3Wrapper.PutObject(ctx, storage.buildKey(ctx, name), data)
}

func (storage s3Storage) buildKey(ctx context.Context, name string) string {
	return fmt.Sprintf("%s/%s", domain.GetChurchID(ctx), name)
}
