resource "aws_lambda_permission" "policy_post_members_create" {
  action        = "lambda:InvokeFunction"
  function_name = var.lambda_arn
  principal     = "apigateway.amazonaws.com"

  source_arn = "arn:aws:execute-api:${var.region}:${var.account_id}:${aws_api_gateway_rest_api.api_gateway.id}/*/POST/members"
}

resource "aws_lambda_permission" "policy_get_members" {
  action        = "lambda:InvokeFunction"
  function_name = var.lambda_arn
  principal     = "apigateway.amazonaws.com"

  source_arn = "arn:aws:execute-api:${var.region}:${var.account_id}:${aws_api_gateway_rest_api.api_gateway.id}/*/GET/members/*"
}

resource "aws_lambda_permission" "policy_put_members" {
  action        = "lambda:InvokeFunction"
  function_name = var.lambda_arn
  principal     = "apigateway.amazonaws.com"

  source_arn = "arn:aws:execute-api:${var.region}:${var.account_id}:${aws_api_gateway_rest_api.api_gateway.id}/*/PUT/members/*/status"
}

resource "aws_lambda_permission" "policy_post_members" {
  action        = "lambda:InvokeFunction"
  function_name = var.lambda_arn
  principal     = "apigateway.amazonaws.com"

  source_arn = "arn:aws:execute-api:${var.region}:${var.account_id}:${aws_api_gateway_rest_api.api_gateway.id}/*/POST/members/search"
}

resource "aws_lambda_permission" "policy_get_reports_members" {
  action        = "lambda:InvokeFunction"
  function_name = var.lambda_arn
  principal     = "apigateway.amazonaws.com"

  source_arn = "arn:aws:execute-api:${var.region}:${var.account_id}:${aws_api_gateway_rest_api.api_gateway.id}/*/GET/reports/members"
}

resource "aws_lambda_permission" "policy_get_reports_members_birthday" {
  action        = "lambda:InvokeFunction"
  function_name = var.lambda_arn
  principal     = "apigateway.amazonaws.com"

  source_arn = "arn:aws:execute-api:${var.region}:${var.account_id}:${aws_api_gateway_rest_api.api_gateway.id}/*/GET/reports/members/birthday"
}

resource "aws_lambda_permission" "policy_get_reports_members_marriage" {
  action        = "lambda:InvokeFunction"
  function_name = var.lambda_arn
  principal     = "apigateway.amazonaws.com"

  source_arn = "arn:aws:execute-api:${var.region}:${var.account_id}:${aws_api_gateway_rest_api.api_gateway.id}/*/GET/reports/members/marriage"
}

resource "aws_lambda_permission" "policy_get_reports_members_legal" {
  action        = "lambda:InvokeFunction"
  function_name = var.lambda_arn
  principal     = "apigateway.amazonaws.com"

  source_arn = "arn:aws:execute-api:${var.region}:${var.account_id}:${aws_api_gateway_rest_api.api_gateway.id}/*/GET/reports/members/legal"
}

resource "aws_lambda_permission" "policy_get_reports_members_classification" {
  action        = "lambda:InvokeFunction"
  function_name = var.lambda_arn
  principal     = "apigateway.amazonaws.com"

  source_arn = "arn:aws:execute-api:${var.region}:${var.account_id}:${aws_api_gateway_rest_api.api_gateway.id}/*/GET/reports/members/classification/*"
}