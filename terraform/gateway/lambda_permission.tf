resource "aws_lambda_permission" "policy-post-members-create" {
  action        = "lambda:InvokeFunction"
  function_name = var.lambda_arn
  principal     = "apigateway.amazonaws.com"

  source_arn = "arn:aws:execute-api:${var.region}:${var.account_id}:${aws_api_gateway_rest_api.api-gateway.id}/*/POST/members"
}

resource "aws_lambda_permission" "policy-get-members" {
  action        = "lambda:InvokeFunction"
  function_name = var.lambda_arn
  principal     = "apigateway.amazonaws.com"

  source_arn = "arn:aws:execute-api:${var.region}:${var.account_id}:${aws_api_gateway_rest_api.api-gateway.id}/*/GET/members/*"
}

resource "aws_lambda_permission" "policy-put-members" {
  action        = "lambda:InvokeFunction"
  function_name = var.lambda_arn
  principal     = "apigateway.amazonaws.com"

  source_arn = "arn:aws:execute-api:${var.region}:${var.account_id}:${aws_api_gateway_rest_api.api-gateway.id}/*/PUT/members/*/status"
}

resource "aws_lambda_permission" "policy-post-members" {
  action        = "lambda:InvokeFunction"
  function_name = var.lambda_arn
  principal     = "apigateway.amazonaws.com"

  source_arn = "arn:aws:execute-api:${var.region}:${var.account_id}:${aws_api_gateway_rest_api.api-gateway.id}/*/POST/members/search"
}

resource "aws_lambda_permission" "policy-get-reports-members" {
  action        = "lambda:InvokeFunction"
  function_name = var.lambda_arn
  principal     = "apigateway.amazonaws.com"

  source_arn = "arn:aws:execute-api:${var.region}:${var.account_id}:${aws_api_gateway_rest_api.api-gateway.id}/*/GET/reports/members"
}

resource "aws_lambda_permission" "policy-get-reports-members-birthday" {
  action        = "lambda:InvokeFunction"
  function_name = var.lambda_arn
  principal     = "apigateway.amazonaws.com"

  source_arn = "arn:aws:execute-api:${var.region}:${var.account_id}:${aws_api_gateway_rest_api.api-gateway.id}/*/GET/reports/members/birthday"
}

resource "aws_lambda_permission" "policy-get-reports-members-marriage" {
  action        = "lambda:InvokeFunction"
  function_name = var.lambda_arn
  principal     = "apigateway.amazonaws.com"

  source_arn = "arn:aws:execute-api:${var.region}:${var.account_id}:${aws_api_gateway_rest_api.api-gateway.id}/*/GET/reports/members/marriage"
}

resource "aws_lambda_permission" "policy-get-reports-members-legal" {
  action        = "lambda:InvokeFunction"
  function_name = var.lambda_arn
  principal     = "apigateway.amazonaws.com"

  source_arn = "arn:aws:execute-api:${var.region}:${var.account_id}:${aws_api_gateway_rest_api.api-gateway.id}/*/GET/reports/members/legal"
}

resource "aws_lambda_permission" "policy-get-reports-members-classification" {
  action        = "lambda:InvokeFunction"
  function_name = var.lambda_arn
  principal     = "apigateway.amazonaws.com"

  source_arn = "arn:aws:execute-api:${var.region}:${var.account_id}:${aws_api_gateway_rest_api.api-gateway.id}/*/GET/reports/members/classification/*"
}