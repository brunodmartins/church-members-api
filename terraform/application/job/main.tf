resource "aws_lambda_function" "lambda_job" {
  function_name = var.lambda_job_name
  role = var.lambda_role_arn
  timeout = 500
  image_uri = var.image_uri
  package_type = "Image"
  environment {
    variables = merge(var.env_var, {"REPORTS_TOPIC" : aws_sns_topic.reports_topic.arn})
  }
}

resource "aws_sns_topic" "reports_topic" {
  name = var.reports_topic_name
}
