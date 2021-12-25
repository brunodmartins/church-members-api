resource "aws_sns_topic" "reports_topic" {
  name = "weekly-report-topic"
}

resource "aws_sns_topic" "alarms_topic" {
  name = "church-members-arlarm-topic"
}

output "reports_topic" {
  value = aws_sns_topic.reports_topic.arn
}

output "alarms_topic" {
  value = aws_sns_topic.alarms_topic.arn
}