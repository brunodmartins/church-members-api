terraform {
  required_providers {
    aws = {
      source = "hashicorp/aws"
    }
  }
}

provider "aws" {
  endpoints {
    s3      = "http://127.0.0.1:4566"
    ec2     = "http://127.0.0.1:4566"
    ecr     = "http://127.0.0.1:4566"
  }
  region                      = "us-east-1"
  skip_credentials_validation = true
  skip_metadata_api_check     = true
  skip_requesting_account_id  = true
  s3_use_path_style = true
}

data "aws_caller_identity" "current" {}
data "aws_region" "current" {}
