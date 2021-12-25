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

variable "region" {
  default = "us-east-1"
}

variable "app_lang" {
  type = string
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

provider "aws" {
  profile = "default"
  region = var.region
}

data "aws_caller_identity" "current" {}
data "aws_region" "current" {}

module "ses" {
  source = "./ses"
  email = var.email_sender
}

module "dynamodb" {
  source = "./dynamodb"
}

module "iam" {
  source = "./iam"
  dynamodb_tables = module.dynamodb.tables_arn
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
  member_history_table_name = module.dynamodb.member_history_table_name
  user_table_name = module.dynamodb.user_table_name
  church_table_name = module.dynamodb.church_table_name
  image_uri = module.ecr.image_id
  lambda_role_arn = module.iam.lambda_role_arn
  topic_arn = module.sns.reports_topic
  app_lang = var.app_lang
  security_token_secret = var.security_token_secret
  security_token_expiration = var.security_token_expiration
  email_sender = var.email_sender
}

module "gateway" {
  source = "./gateway"
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

output "gateway_id" {
  value = module.gateway.gateway_id
}