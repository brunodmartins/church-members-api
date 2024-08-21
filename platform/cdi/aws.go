package cdi

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sesv2"
	"github.com/brunodmartins/church-members-api/platform/aws/wrapper"
	"github.com/spf13/viper"
)

func provideDynamoDB() *dynamodb.Client {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic("unable to load SDK config, " + err.Error())
	}
	if viper.Get("cloud") == "LOCAL" {
		return dynamodb.NewFromConfig(cfg, func(options *dynamodb.Options) {
			options.BaseEndpoint = aws.String("http://127.0.0.1:4566")
		})
	}
	return dynamodb.NewFromConfig(cfg)
}

func provideSNS() *sns.Client {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic("unable to load SDK config, " + err.Error())
	}
	return sns.NewFromConfig(cfg)
}

func provideSES() wrapper.SESAPI {
	if viper.Get("cloud") == "LOCAL" {
		return wrapper.NewMockSESAPI()
	}
	mySession := session.Must(session.NewSession())
	return sesv2.New(mySession)
}

func provideS3() *s3.Client {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic("unable to load SDK config, " + err.Error())
	}
	if viper.Get("cloud") == "LOCAL" {
		return s3.NewFromConfig(cfg, func(options *s3.Options) {
			options.BaseEndpoint = aws.String("http://127.0.0.1:4566")
		})
	}
	return s3.NewFromConfig(cfg)
}
