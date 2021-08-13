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

variable "topic_arn" {
  type = string
}

resource "aws_lambda_function" "lambda_api" {
  function_name = "church-members-api-lambda"
  role          = var.lambda_role_arn
  timeout       = 500
  image_uri     = var.image_uri
  package_type  = "Image"
  environment {
    variables = {
      "SERVER" : "AWS",
      "APPLICATION" : "API"
      "CHURCH_NAME" : "",
      "CHURCH_NAME_SHORT" : "",
      "APP_LANG" : "pt-BR",
      "TABLE_MEMBER" : var.member_table_name,
      "TABLE_MEMBER_HISTORY" : var.member_history_table_name,
      "JOBS_DAILY_PHONE" : "",
      "REPORTS_TOPIC" : var.topic_arn,
    }
  }
}

resource "aws_lambda_function" "lambda_job" {
  function_name = "church-members-job-lambda"
  role          = var.lambda_role_arn
  timeout       = 500
  image_uri     = var.image_uri
  package_type  = "Image"
  environment {
    variables = {
      "SERVER" : "AWS",
      "APPLICATION" : "JOB"
      "CHURCH_NAME" : "",
      "CHURCH_NAME_SHORT" : "",
      "APP_LANG" : "pt-BR",
      "TABLE_MEMBER" : var.member_table_name,
      "TABLE_MEMBER_HISTORY" : var.member_history_table_name,
      "JOBS_DAILY_PHONE" : "",
      "REPORTS_TOPIC" : var.topic_arn,
    }
  }
}


output "lambda_arn" {
  value = aws_lambda_function.lambda_api.arn
}

output "lambda_name" {
  value = aws_lambda_function.lambda_api.function_name
}

output "lambda_job_arn" {
  value = aws_lambda_function.lambda_job.arn
}