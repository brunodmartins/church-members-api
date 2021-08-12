variable "region" {
  type = string
}

variable "account_id" {
  type = string
}

variable "cognito_user_pool_arn" {
  type = string
}

variable "lambda_name" {
  type = string
}

variable "lambda_arn" {
  type = string
}

data "template_file" "aws_api_swagger" {
  template = file("${path.module}/swagger-terraform.json")
  vars = {
    aws_region     = var.region
    aws_account_id = var.account_id
    lambda_id      = var.lambda_name
    cognito_pool   = var.cognito_user_pool_arn
  }
}

resource "aws_api_gateway_rest_api" "api_gateway" {
  name        = "church-members-api-gw"
  description = "church-members-api API gateway"
  body        = data.template_file.aws_api_swagger.rendered

  endpoint_configuration {
    types = ["REGIONAL"]
  }
}

resource "aws_api_gateway_deployment" "api_deployment" {
  rest_api_id = aws_api_gateway_rest_api.api_gateway.id

  triggers = {
    redeployment = sha1(jsonencode(aws_api_gateway_rest_api.api_gateway.body))
  }

  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_api_gateway_stage" "api_stage" {
  depends_on = [
    aws_api_gateway_deployment.api_deployment
  ]
  deployment_id = aws_api_gateway_deployment.api_deployment.id
  rest_api_id   = aws_api_gateway_rest_api.api_gateway.id
  stage_name    = "prod"
}

resource "aws_api_gateway_authorizer" "authorizer" {
  name        = "authorizer"
  rest_api_id = aws_api_gateway_rest_api.api_gateway.id
  type        = "COGNITO_USER_POOLS"
  provider_arns = [
    var.cognito_user_pool_arn
  ]
}


output "gateway_id" {
  value = aws_api_gateway_rest_api.api_gateway.id
}