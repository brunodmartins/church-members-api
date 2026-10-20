# When the job module was split, these production resources moved to job_schedule.
# Keep these mappings so Terraform preserves existing resources in state instead
# of destroying and recreating them during the module refactor.
moved {
  from = module.job.aws_cloudwatch_metric_alarm.job-error
  to   = module.job_schedule.aws_cloudwatch_metric_alarm.job-error
}

moved {
  from = module.job.aws_sns_topic.alarms_topic
  to   = module.job_schedule.aws_sns_topic.alarms_topic
}

moved {
  from = module.job.aws_cloudwatch_event_rule.weekly_report
  to   = module.job_schedule.aws_cloudwatch_event_rule.weekly_report
}

moved {
  from = module.job.aws_cloudwatch_event_rule.daily_report
  to   = module.job_schedule.aws_cloudwatch_event_rule.daily_report
}

moved {
  from = module.job.aws_cloudwatch_event_rule.daily_members_novelties
  to   = module.job_schedule.aws_cloudwatch_event_rule.daily_members_novelties
}

moved {
  from = module.job.aws_cloudwatch_event_target.weekly_lambda
  to   = module.job_schedule.aws_cloudwatch_event_target.weekly_lambda
}

moved {
  from = module.job.aws_cloudwatch_event_target.daily_lambda
  to   = module.job_schedule.aws_cloudwatch_event_target.daily_lambda
}

moved {
  from = module.job.aws_cloudwatch_event_target.daily_members_novelties_lambda
  to   = module.job_schedule.aws_cloudwatch_event_target.daily_members_novelties_lambda
}

moved {
  from = module.job.aws_lambda_permission.cloudwatch-daily-members-novelties
  to   = module.job_schedule.aws_lambda_permission.cloudwatch-daily-members-novelties
}

moved {
  from = module.job.aws_lambda_permission.cloudwatch-daily
  to   = module.job_schedule.aws_lambda_permission.cloudwatch-daily
}

moved {
  from = module.job.aws_lambda_permission.cloudwatch-weekly
  to   = module.job_schedule.aws_lambda_permission.cloudwatch-weekly
}
