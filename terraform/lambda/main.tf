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

variable "user_table_name" {
  type = string
}

variable "church_table_name" {
  type = string
}

variable "topic_arn" {
  type = string
}

variable "app_lang" {
  type = string
}

variable "church_name" {
  type = string
}

variable "church_name_short" {
  type = string
}

variable "security_token_secret" {
  type = string
}

variable "security_token_expiration" {
  type = string
}

variable "email_sender" {
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
      "APP_LANG" : var.app_lang,
      "EMAIL_SENDER": var.email_sender,
      "TABLE_MEMBER" : var.member_table_name,
      "TABLE_MEMBER_HISTORY" : var.member_history_table_name,
      "TABLE_USER" : var.user_table_name,
      "TABLE_CHURCH": var.church_table_name,
      "REPORTS_TOPIC" : var.topic_arn,
      "TOKEN_SECRET" : var.security_token_secret,
      "TOKEN_EXPIRATION" : var.security_token_expiration,
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
      "APP_LANG" : var.app_lang,
      "EMAIL_SENDER": var.email_sender,
      "TABLE_MEMBER" : var.member_table_name,
      "TABLE_USER" : var.user_table_name,
      "TABLE_MEMBER_HISTORY" : var.member_history_table_name,
      "TABLE_CHURCH": var.church_table_name,
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