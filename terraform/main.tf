terraform {
  required_providers {
    aws = {
      source = "hashicorp/aws"
    }
  }

  backend "remote" {
    hostname = "app.terraform.io"
    organization = "church-members-api"

    workspaces {
      name = "church-members-api"
    }
  }
}

variable "security_token_secret" {
  type = string
}

variable "security_token_expiration" {
  type = string
}

variable "email_sender" {
  type = string
}

provider "aws" {}

data "aws_caller_identity" "current" {}
data "aws_region" "current" {}

module "ses" {
  source = "./ses"
  email = var.email_sender
}

module "dynamodb" {
  source = "./dynamodb"
  member_table_name = "member_v2"
  user_table_name = "user_v2"
  church_table_name = "church"
}

module "iam" {
  source = "./iam"
  dynamodb_tables = module.dynamodb.tables_arn
  bucket_arn = module.s3.bucket_arn
  role_name = "church-members-api"
}

module "ecr" {
  source = "./ecr"
}

module "sns" {
  source = "./sns"
}

module "lambda" {
  source = "./lambda"
  member_table_name = module.dynamodb.member_table_name
  user_table_name = module.dynamodb.user_table_name
  church_table_name = module.dynamodb.church_table_name
  image_uri = module.ecr.image_id
  lambda_role_arn = module.iam.lambda_role_arn
  topic_arn = module.sns.reports_topic
  security_token_secret = var.security_token_secret
  security_token_expiration = var.security_token_expiration
  email_sender = var.email_sender
  bucket_name = module.s3.bucket_name
  lambda_api_name = "church-members-api-lambda"
  lambda_job_name = "church-members-job-lambda"
}

module "gateway" {
  source = "./gateway"
  gateway_name = "church-members-api-gw"
  region = data.aws_region.current.name
  account_id = data.aws_caller_identity.current.account_id
  lambda_name = module.lambda.lambda_name
  lambda_arn = module.lambda.lambda_arn
}

module "eventbridge" {
  source = "./eventbridge"
  lambda_arn = module.lambda.lambda_job_arn
}

module "cloudwatch" {
  source = "./cloudwatch"
  job_function = module.lambda.lambda_job_name
  sns_topic = module.sns.alarms_topic
}

module "s3" {
  source = "./s3"
  bucket_name = "church-members-data"
}

output "gateway_id" {
  value = module.gateway.gateway_id
}