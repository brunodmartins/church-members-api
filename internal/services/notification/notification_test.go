package notification

import (
	"errors"
	mock_wrapper "github.com/BrunoDM2943/church-members-api/platform/aws/wrapper/mock"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSnsService_NotifyTopic(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	snsMock := mock_wrapper.NewMockSNSAPI(ctrl)
	const topicARN = "my-topic-arn"
	const message = "my-message"

	service := NewService(snsMock, topicARN)

	t.Run("Success", func(t *testing.T) {
		input := &sns.PublishInput{
			Message:  aws.String(message),
			TopicArn: aws.String(topicARN),
		}
		snsMock.EXPECT().Publish(gomock.Any(), gomock.Eq(input), gomock.Any()).Return(nil, nil)
		assert.Nil(t, service.NotifyTopic(message))
	})
	t.Run("Fail", func(t *testing.T) {
		input := &sns.PublishInput{
			Message:  aws.String(message),
			TopicArn: aws.String(topicARN),
		}
		snsMock.EXPECT().Publish(gomock.Any(), gomock.Eq(input), gomock.Any()).Return(nil, errors.New("generic error"))
		assert.NotNil(t, service.NotifyTopic(message))
	})
}

func TestSnsService_NotifyMobile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	snsMock := mock_wrapper.NewMockSNSAPI(ctrl)
	const topicARN = "my-topic-arn"
	const phone = "my-phone"
	const message = "my-message"

	service := NewService(snsMock, topicARN)

	t.Run("Success", func(t *testing.T) {
		input := &sns.PublishInput{
			Message:     aws.String(message),
			PhoneNumber: aws.String(phone),
		}
		snsMock.EXPECT().Publish(gomock.Any(), gomock.Eq(input), gomock.Any()).Return(nil, nil)
		assert.Nil(t, service.NotifyMobile(message, phone))
	})
	t.Run("Fail", func(t *testing.T) {
		input := &sns.PublishInput{
			Message:     aws.String(message),
			PhoneNumber: aws.String(phone),
		}
		snsMock.EXPECT().Publish(gomock.Any(), gomock.Eq(input), gomock.Any()).Return(nil, errors.New("generic error"))
		assert.NotNil(t, service.NotifyMobile(message, phone))
	})
}
