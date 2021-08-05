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

provider "aws" {
  profile = "default"
  region  = var.region
}

data "aws_caller_identity" "current" {}
data "aws_region" "current" {}

locals {
  prefix              = "git"
  app_name            = "church-members-api"
  account_id          = data.aws_caller_identity.current.account_id
  ecr_repository_name = "${local.app_name}-container"
  ecr_image_tag       = "latest"
  swagger-file-path   = "swagger-terraform.json"
}

resource "aws_dynamodb_table" "member-table" {
  name           = "member"
  billing_mode   = "PROVISIONED"
  read_capacity  = 5
  write_capacity = 5
  hash_key       = "id"

  attribute {
    name = "id"
    type = "S"
  }
}

resource "aws_dynamodb_table" "member-history-table" {
  name           = "member_history"
  billing_mode   = "PROVISIONED"
  read_capacity  = 5
  write_capacity = 5
  hash_key       = "id"

  attribute {
    name = "id"
    type = "S"
  }
}

resource "aws_iam_policy" "church-members-api-policy" {
  name        = "${local.app_name}-policy"
  description = "This policy allow ${local.app_name} full execution"
  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "dynamodb:PutItem",
          "dynamodb:UpdateItem",
          "dynamodb:DeleteItem",
          "dynamodb:GetItem",
          "dynamodb:Scan",
        ]
        Resource = [
          aws_dynamodb_table.member-table.arn,
          aws_dynamodb_table.member-history-table.arn,
        ]
      },
      {
        Effect = "Allow"
        Action = [
          "logs:CreateLogGroup",
          "logs:CreateLogStream",
          "logs:PutLogEvents",
        ]
        Resource = "*"
      }
    ]
  })
}

resource "aws_iam_role" "church-members-api-role" {
  name = "${local.app_name}-role"
  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Principal = {
          Service = "lambda.amazonaws.com"
        }
      }
  ] })
}

resource "aws_iam_role_policy_attachment" "attach-policy" {
  role       = aws_iam_role.church-members-api-role.name
  policy_arn = aws_iam_policy.church-members-api-policy.arn
}

resource "aws_ecr_repository" "repo" {
  name = local.ecr_repository_name
}

data "archive_file" "internal" {
  type        = "zip"
  source_dir  = "internal/"
  output_path = "internal.zip"
}

data "archive_file" "cmd" {
  type        = "zip"
  source_dir  = "internal/"
  output_path = "cmd.zip"
}


data "aws_ecr_image" "lambda_image" {
  repository_name = local.ecr_repository_name
  image_tag       = local.ecr_image_tag
}

resource "aws_lambda_function" "lambda" {
  depends_on = [
    aws_dynamodb_table.member-table,
    aws_dynamodb_table.member-history-table
  ]
  function_name = "${local.app_name}-lambda"
  role          = aws_iam_role.church-members-api-role.arn
  timeout       = 500
  image_uri     = "${aws_ecr_repository.repo.repository_url}@${data.aws_ecr_image.lambda_image.id}"
  package_type  = "Image"
  environment {
    variables = {
      "SERVER" : "AWS",
      "SCOPE" : "prod"
      "APP_LANG" : "pt-BR",
      "CHURCH_NAME" : "",
      "TABLES_MEMBER": aws_dynamodb_table.member-table.name,
      "TABLES_MEMBER_HISTORY": aws_dynamodb_table.member-history-table.name,
    }
  }
}

data "template_file" "aws_api_swagger" {
  template = file(local.swagger-file-path)
  vars = {
    aws_region     = data.aws_region.current.name
    aws_account_id = data.aws_caller_identity.current.account_id
    lambda_id      = "${local.app_name}-lambda"
    cognito_pool   = aws_cognito_user_pool.user-pool.arn
  }
}

resource "aws_api_gateway_rest_api" "api-gateway" {
  name        = "${local.app_name}-gw"
  description = "${local.app_name} API gateway"
  body        = data.template_file.aws_api_swagger.rendered

  endpoint_configuration {
    types = ["REGIONAL"]
  }
}

resource "aws_api_gateway_deployment" "api-deployment" {
  rest_api_id = aws_api_gateway_rest_api.api-gateway.id

  triggers = {
    redeployment = sha1(jsonencode(aws_api_gateway_rest_api.api-gateway.body))
  }

  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_api_gateway_stage" "api-stage" {
  depends_on = [
    aws_api_gateway_deployment.api-deployment
  ]
  deployment_id = aws_api_gateway_deployment.api-deployment.id
  rest_api_id   = aws_api_gateway_rest_api.api-gateway.id
  stage_name    = "prod"
}

resource "aws_lambda_permission" "policy-post-members-create" {
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.lambda.arn
  principal     = "apigateway.amazonaws.com"

  source_arn = "arn:aws:execute-api:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:${aws_api_gateway_rest_api.api-gateway.id}/*/POST/members"
}

resource "aws_lambda_permission" "policy-get-members" {
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.lambda.arn
  principal     = "apigateway.amazonaws.com"

  source_arn = "arn:aws:execute-api:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:${aws_api_gateway_rest_api.api-gateway.id}/*/GET/members/*"
}

resource "aws_lambda_permission" "policy-put-members" {
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.lambda.arn
  principal     = "apigateway.amazonaws.com"

  source_arn = "arn:aws:execute-api:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:${aws_api_gateway_rest_api.api-gateway.id}/*/PUT/members/*/status"
}

resource "aws_lambda_permission" "policy-post-members" {
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.lambda.arn
  principal     = "apigateway.amazonaws.com"

  source_arn = "arn:aws:execute-api:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:${aws_api_gateway_rest_api.api-gateway.id}/*/POST/members/search"
}

resource "aws_lambda_permission" "policy-get-reports-members" {
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.lambda.arn
  principal     = "apigateway.amazonaws.com"

  source_arn = "arn:aws:execute-api:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:${aws_api_gateway_rest_api.api-gateway.id}/*/GET/reports/members"
}

resource "aws_lambda_permission" "policy-get-reports-members-birthday" {
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.lambda.arn
  principal     = "apigateway.amazonaws.com"

  source_arn = "arn:aws:execute-api:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:${aws_api_gateway_rest_api.api-gateway.id}/*/GET/reports/members/birthday"
}

resource "aws_lambda_permission" "policy-get-reports-members-marriage" {
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.lambda.arn
  principal     = "apigateway.amazonaws.com"

  source_arn = "arn:aws:execute-api:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:${aws_api_gateway_rest_api.api-gateway.id}/*/GET/reports/members/marriage"
}

resource "aws_lambda_permission" "policy-get-reports-members-legal" {
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.lambda.arn
  principal     = "apigateway.amazonaws.com"

  source_arn = "arn:aws:execute-api:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:${aws_api_gateway_rest_api.api-gateway.id}/*/GET/reports/members/legal"
}

resource "aws_lambda_permission" "policy-get-reports-members-classification" {
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.lambda.arn
  principal     = "apigateway.amazonaws.com"

  source_arn = "arn:aws:execute-api:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:${aws_api_gateway_rest_api.api-gateway.id}/*/GET/reports/members/classification/*"
}

resource "aws_cognito_user_pool" "user-pool" {
  name = "church-members-user-pool"
  admin_create_user_config {
    allow_admin_create_user_only = true
  }
}

resource "aws_api_gateway_authorizer" "authorizer" {
  depends_on = [
    aws_api_gateway_rest_api.api-gateway
  ]
  name          = "authorizer"
  rest_api_id   = aws_api_gateway_rest_api.api-gateway.id
  type          = "COGNITO_USER_POOLS"
  provider_arns = [
    aws_cognito_user_pool.user-pool.arn
  ]
}

output "api-gateway" {
  value = aws_api_gateway_rest_api.api-gateway.id
}