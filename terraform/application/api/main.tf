data "aws_caller_identity" "current" {}
data "aws_region" "current" {}


resource "aws_lambda_function" "lambda_api" {
  function_name = var.lambda_api_name
  role = var.lambda_role_arn
  timeout = 500
  image_uri = var.image_uri
  package_type = "Image"
  environment {
    variables = var.env_var
  }
}

data "template_file" "aws_api_swagger" {
  template = file("${path.module}/swagger-terraform.json")
  vars = {
    aws_region = data.aws_region.current.name
    aws_account_id = data.aws_caller_identity.current.account_id
    lambda_id = aws_lambda_function.lambda_api.function_name
  }
}

resource "aws_api_gateway_rest_api" "api_gateway" {
  name = var.gateway_name
  description = "church-members-api API gateway"
  body = data.template_file.aws_api_swagger.rendered
  minimum_compression_size = 1
  binary_media_types = [
    "*/*",
    "application/pdf"
  ]

  endpoint_configuration {
    types = [
      "REGIONAL"]
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
  rest_api_id = aws_api_gateway_rest_api.api_gateway.id
  cache_cluster_size = "0.5"
  stage_name = "prod"
}

resource "aws_lambda_permission" "policy_any_proxy" {
  action        = "lambda:InvokeFunction"
  function_name =  aws_lambda_function.lambda_api.arn
  principal     = "apigateway.amazonaws.com"
  source_arn = "arn:aws:execute-api:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:${aws_api_gateway_rest_api.api_gateway.id}/*/*/*"
}