terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 6.60.0"
    }
  }

  backend "remote" {
    hostname     = "app.terraform.io"
    organization = "church-members-api"

    workspaces {
      name = "production"
    }
  }
}

provider "aws" {
  default_tags {
    tags = {
      Environment = "Production"
      Application = "church-members-api"
    }
  }
}

data "aws_caller_identity" "current" {}
data "aws_region" "current" {}
