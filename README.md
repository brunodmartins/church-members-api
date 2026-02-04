# church-members-api ⛪️

![Go](https://github.com/brunodmartins/church-members-api/workflows/Go/badge.svg)
[![codecov](https://codecov.io/gh/brunodmartins/church-members-api/branch/master/graph/badge.svg)](https://codecov.io/gh/brunodmartins/church-members-api)  [![Go Report Card](https://goreportcard.com/badge/github.com/brunodmartins/church-members-api?style=flat-square)](https://goreportcard.com/report/github.com/brunodmartins/church-members-api)
![Deploy](https://github.com/brunodmartins/church-members-api/workflows/Docker%20Image%20CI/badge.svg)

A simple application to manage a church's members.

## Features 💻

- Member's control (register & search)
- Member's report
  - Legal members (all except children)
  - All members (including children)
  - By classification
  - Birthdays list
  - Marriage list
- Notification jobs
    - Weekly births through email
    - Dailly birthdays through SMS
- i18n
  - Supports pt-BR and en-US

## Libraries ⚙️

- [Viper](https://github.com/spf13/viper)
- [go-age](https://github.com/bearbin/go-age)
- [AWS SDK v2](https://github.com/aws/aws-sdk-go-v2)
- [gopdf](https://github.com/signintech/gopdf)
- [nicksnyder-i18n](https://github.com/nicksnyder/go-i18n/v2/i18n)
- [logrus](https://github.com/sirupsen/logrus)

## Deploy ✈️

The application was build upon a Docker image, but rely mostly on AWS resources to work. To deploy on AWS, simple configure a **Terraform**  with some vars described bellow and run the script and it's all done!

### AWS ☁️

The following resources are used on AWS

- **DynamoDB** to store churchs, users and members information
- **S3** to store report files
- **Event Bridge** to CRON the jobs
- **SNS** to send notifications from jobs
- **ECR** to store a private image of the application
- **Lambda** to run the serverless application (both API and JOB)
- **API Gateway** to provide a RESTfull interface and authorize access

![GitHub Logo](/docs/architecture.png)

### Configuration 🛠

The following configuration are required through **Terraform vars**

|Terraform var|Environment variable|Description
|-|-|-|
|Hard coded on Terraform|SERVER|Used to define the environment where the application will run. Defaults to **AWS**|
|Hard coded on Terraform|APPLICATION|Used to defined the Lambda type: **API** or **JOB**|
|Terraform takes it from dynamo resource|TABLE_MEMBER|DynamoDB members table|
|Terraform takes it from dynamo resource|TABLE_USER|DynamoDB users table|
|Terraform takes it from dynamo resource|TABLE_CHURCH|DynamoDB churchs table|
|Terraform takes it from dynamo resource|TABLE_PARTICIPANT|DynamoDB participants table|
|Terraform takes it from SNS resource|REPORTS_TOPIC|The topic to notify weekly jobs through email|
|Terraform takes it from S3 resource|STORAGE|The S3 bucket|

## Support ✉️

You can create a PR for it =)
