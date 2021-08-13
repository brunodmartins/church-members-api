package notification

import (
	"context"
	"github.com/BrunoDM2943/church-members-api/platform/aws/wrapper"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/sirupsen/logrus"
)

//go:generate mockgen -source=./notification.go -destination=./mock/notification_mock.go
//Service provide operations for notification
type Service interface {
	//NotifyTopic send a notification to a defined topic
	NotifyTopic(text string) error
	//NotifyMobile send a notification to a mobile phone
	NotifyMobile(text string, phone string) error
}

type snsService struct {
	api   wrapper.SNSAPI
	topic string
}

//NewService builds a new notification service
func NewService(api wrapper.SNSAPI, topic string) Service {
	return &snsService{
		api:   api,
		topic: topic,
	}
}

func (service *snsService) NotifyTopic(text string) error {
	input := &sns.PublishInput{
		Message:  aws.String(text),
		TopicArn: aws.String(service.topic),
	}
	logrus.Info("Send notification to topic")
	_, err := service.api.Publish(context.TODO(), input)
	if err != nil {
		logrus.Info("Notification sent!")
	}
	return err
}

func (service *snsService) NotifyMobile(text string, phone string) error {
	input := &sns.PublishInput{
		Message:     aws.String(text),
		PhoneNumber: aws.String(phone),
	}
	logrus.Info("Send notification to phone")
	_, err := service.api.Publish(context.TODO(), input)
	if err != nil {
		logrus.Info("Notification sent!")
	}
	return err
}
