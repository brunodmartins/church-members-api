# church-members-api ⛪️

![Go](https://github.com/BrunoDM2943/church-members-api/workflows/Go/badge.svg)
[![codecov](https://codecov.io/gh/BrunoDM2943/church-members-api/branch/master/graph/badge.svg)](https://codecov.io/gh/BrunoDM2943/church-members-api)  [![Go Report Card](https://goreportcard.com/badge/github.com/BrunoDM2943/church-members-api?style=flat-square)](https://goreportcard.com/report/github.com/BrunoDM2943/church-members-api)
![Deploy](https://github.com/BrunoDM2943/church-members-api/workflows/Docker%20Image%20CI/badge.svg)


A simple application to manage a church's members. 

## Features 💻

- Member's control (register & search)
- Member's report
  - Legal members (all except children)
  - All members (including children)
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

- **DynamoDB** for store members information
- **Event Bridge** to CRON the jobs
- **SNS** to send notifications from jobs
- **ECR** to store a private image of the application
- **Lambda** to run the serverless application (both API and JOB)
- **API Gateway** to provide a RESTfull interface and authorize access
- **Cognito** to authenticate users

![GitHub Logo](/docs/architecture.png)

### Configuration 🛠

The following configuration are required through **Terraform vars**

|Terraform var|Environment variable|Description
|-|-|-|
|Hard coded on Terraform|SERVER|Used to define the environment where the application will run. Defaults to **AWS**|
|Hard coded on Terraform|APPLICATION|Used to defined the Lambda type: **API** or **JOB**|
|church_name|CHURCH_NAME|The church name|
|church_name_short|CHURCH_NAME_SHORT|The abbreviation of the church name|
|app_lang|APP_LANG|The application language. See support languages above on features|
|Terraform takes it from dynamo resource|TABLE_MEMBER|DynamoDB members table|
|Terraform takes it from dynamo resource|TABLE_MEMBER_HISTORY|DynamoDB members history table|
|-|JOBS_DAILY_PHONE|A list of phone numbers separated by comma to receive the SMS notification|
|Terraform takes it from SNS resource|REPORTS_TOPIC|The topic to notify weekly jobs through email|

## Support ✉️

You can create a PR for it =) 