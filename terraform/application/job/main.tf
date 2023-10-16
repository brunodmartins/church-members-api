data "aws_caller_identity" "current" {}
data "aws_region" "current" {}


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
  name = "weekly-report-topic"
}

resource "aws_cloudwatch_metric_alarm" "job-error" {
  alarm_name                = "jobs-execution-error"
  comparison_operator       = "GreaterThanOrEqualToThreshold"
  evaluation_periods        = "1"
  metric_name               = "Errors"
  namespace                 = "AWS/Lambda"
  period                    = "3600"
  statistic                 = "Sum"
  threshold                 = "1"
  alarm_description         = "This alarm triggers when a single lambda execution fails"
  treat_missing_data        = "notBreaching"
  alarm_actions             = [aws_sns_topic.alarms_topic.arn]
  dimensions = {
    FunctionName = aws_lambda_function.lambda_job.function_name
  }
}

resource "aws_sns_topic" "alarms_topic" {
  name = "church-members-arlarm-topic"
}

resource "aws_cloudwatch_event_rule" "weekly_report" {
  name                = "weekly-birthdays-report"
  description         = "Weekly sunday birthdays/marrigage report"
  schedule_expression = "cron(0 13 ? * SUN *)"
}

resource "aws_cloudwatch_event_rule" "daily_report" {
  name                = "daily-birthdays-report"
  description         = "Daily birthdays report"
  schedule_expression = "cron(0 13 ? * * *)"
}


resource "aws_cloudwatch_event_target" "weekly_lambda" {
  arn   = aws_lambda_function.lambda_job.arn
  rule  = aws_cloudwatch_event_rule.weekly_report.name
  input = <<DOC
    {
        "name": "WEEKLY_BIRTHDAYS"
    }
    DOC
}

resource "aws_cloudwatch_event_target" "daily_lambda" {
  arn   = aws_lambda_function.lambda_job.arn
  rule  = aws_cloudwatch_event_rule.daily_report.name
  input = <<DOC
    {
        "name": "DAILY_BIRTHDAYS"
    }
    DOC
}

resource "aws_lambda_permission" "cloudwatch-daily" {
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.lambda_job.arn
  principal     = "events.amazonaws.com"
  source_arn    = aws_cloudwatch_event_rule.daily_report.arn
}

resource "aws_lambda_permission" "cloudwatch-weekly" {
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.lambda_job.arn
  principal     = "events.amazonaws.com"
  source_arn    = aws_cloudwatch_event_rule.weekly_report.arn
}
