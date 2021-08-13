resource "aws_sns_topic" "reports_topic" {
  name = "weekly-report-topic"
}

output "topic_arn" {
  value = aws_sns_topic.reports_topic.arn
}