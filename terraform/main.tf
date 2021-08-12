terraform {
  required_providers {
    aws = {
      source = "hashicorp/aws"
    }
  }

  backend "remote" {
    hostname     = "app.terraform.io"
    organization = "church-members-api"

    workspaces {
      name = "church-members-api"
    }
  }
}

variable "region" {
  default = "us-east-1"
}

provider "aws" {
  profile = "default"
  region  = var.region
}

data "aws_caller_identity" "current" {}
data "aws_region" "current" {}


module "dynamodb" {
  source = "./dynamodb"
}

module "iam" {
  source          = "./iam"
  dynamodb_tables = module.dynamodb.tables_arn
}

module "ecr" {
  source = "./ecr"
}

module "lambda" {
  source                     = "./lambda"
  member_table_name         = module.dynamodb.member_table_name
  member_history_table_name = module.dynamodb.member_history_table_name
  image_uri                  = module.ecr.image_id
  lambda_role_arn            = module.iam.lambda_role_arn
}

module "cognito" {
  source = "./cognito"
}

module "gateway" {
  source                = "./gateway"
  region                = data.aws_region.current.name
  account_id            = data.aws_caller_identity.current.account_id
  cognito_user_pool_arn = module.cognito.user_pool_arn
  lambda_name           = module.lambda.lambda_name
  lambda_arn            = module.lambda.lambda_arn
}

output "gateway_id" {
  value = module.gateway.gateway_id
}