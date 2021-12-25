variable "job_function" {
  type = string
}

variable "sns_topic" {
  type = string
}

resource "aws_cloudwatch_metric_alarm" "job-error" {
  alarm_name                = "jobs-execution-error"
  comparison_operator       = "GreaterThanOrEqualToThreshold"
  evaluation_periods        = "1"
  metric_name               = "Errors"
  namespace                 = "AWS/Lambda"
  period                    = "3600"
  statistic                 = "SampleCount"
  threshold                 = "1"
  alarm_description         = "This alarm triggers when a single lambda execution fails"
  treat_missing_data        = "notBreaching"
  alarm_actions             = [var.sns_topic]
  dimensions = {
    FunctionName = var.job_function
  }
}
