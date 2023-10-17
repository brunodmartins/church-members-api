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
      name = "beta"
    }
  }
}

provider "aws" {
  default_tags {
    tags = {
      Environment = "Beta"
      Application = "church-members-api"
    }
  }
}

data "aws_caller_identity" "current" {}
data "aws_region" "current" {}
