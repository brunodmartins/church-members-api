variable "lambda_role_arn" {
  type = string
}

variable "image_uri" {
  type = string
}

variable "member_table_name" {
  type = string
}

variable "member_history_table_name" {
  type = string
}

resource "aws_lambda_function" "lambda" {
  function_name = "church-members-api-lambda"
  role          = var.lambda_role_arn
  timeout       = 500
  image_uri     = var.image_uri
  package_type  = "Image"
  environment {
    variables = {
      "SERVER" : "AWS",
      "SCOPE" : "prod"
      "APP_LANG" : "pt-BR",
      "CHURCH_NAME" : "",
      "TABLES_MEMBER" : var.member_table_name,
      "TABLES_MEMBER_HISTORY" : var.member_history_table_name,
    }
  }
}

output "lambda_arn" {
  value = aws_lambda_function.lambda.arn
}

output "lambda_name" {
  value = aws_lambda_function.lambda.function_name
}