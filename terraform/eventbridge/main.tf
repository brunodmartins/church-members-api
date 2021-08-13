variable "lambda_arn" {
  type = string
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
  arn   = var.lambda_arn
  rule  = aws_cloudwatch_event_rule.weekly_report.name
  input = <<DOC
    {
        "name": "WEEKLY_BIRTHDAYS"
    }
    DOC
}

resource "aws_cloudwatch_event_target" "daily_lambda" {
  arn   = var.lambda_arn
  rule  = aws_cloudwatch_event_rule.daily_report.name
  input = <<DOC
    {
        "name": "DAILY_BIRTHDAYS"
    }
    DOC
}

resource "aws_lambda_permission" "cloudwatch-daily" {
  action        = "lambda:InvokeFunction"
  function_name = var.lambda_arn
  principal     = "events.amazonaws.com"
  source_arn    = aws_cloudwatch_event_rule.daily_report.arn
}

resource "aws_lambda_permission" "cloudwatch-weekly" {
  action        = "lambda:InvokeFunction"
  function_name = var.lambda_arn
  principal     = "events.amazonaws.com"
  source_arn    = aws_cloudwatch_event_rule.weekly_report.arn
}