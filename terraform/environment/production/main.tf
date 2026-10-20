locals {
  ecr_image_tag       = trimspace(file("${path.root}/../../../.version"))
  ecr_repository_name = "church-members-api-container"
}

data "aws_ecr_repository" "repository" {
  name = local.ecr_repository_name
}

data "aws_ecr_image" "lambda_image" {
  repository_name = local.ecr_repository_name
  image_tag       = local.ecr_image_tag
}


module "dynamodb_tables" {
  source                 = "../../storage/dynamodb"
  member_table_name      = var.member_table_name
  user_table_name        = var.user_table_name
  church_table_name      = var.church_table_name
  participant_table_name = var.participant_table_name
}

module "iam_roles" {
  source          = "../../security/iam"
  bucket_arn      = module.s3_bucket.bucket_arn
  dynamodb_tables = module.dynamodb_tables.tables_arn
  queue_arn       = module.sqs_queue.queue_arn
  role_name       = var.role_name
}

module "s3_bucket" {
  source      = "../../storage/s3"
  bucket_name = var.bucket_name
}

module "sqs_queue" {
  source     = "../../infra/sqs"
  queue_name = var.queue_name
}


module "api" {
  source          = "../../application/api"
  lambda_api_name = var.lambda_api_name
  lambda_role_arn = module.iam_roles.lambda_role_arn
  image_uri       = "${data.aws_ecr_repository.repository.repository_url}@${data.aws_ecr_image.lambda_image.id}"
  env_var = {
    "SERVER" : "AWS",
    "APPLICATION" : "API"
    "EMAIL_SENDER" : var.email_sender,
    "TABLE_MEMBER" : module.dynamodb_tables.member_table_name,
    "TABLE_USER" : module.dynamodb_tables.user_table_name,
    "TABLE_CHURCH" : module.dynamodb_tables.church_table_name,
    "TABLE_PARTICIPANT" : module.dynamodb_tables.participant_table_name,
    "TOKEN_SECRET" : var.security_token_secret,
    "TOKEN_EXPIRATION" : var.security_token_expiration,
    "STORAGE" : module.s3_bucket.bucket_name,
    "QUEUE_URL" : module.sqs_queue.queue_url,
    "API_ENDPOINT" : var.gateway_url
  }
  gateway_name = var.gateway_name
}

module "job" {
  source          = "../../application/job"
  lambda_job_name = var.lambda_job_name
  lambda_role_arn = module.iam_roles.lambda_role_arn
  image_uri       = "${data.aws_ecr_repository.repository.repository_url}@${data.aws_ecr_image.lambda_image.id}"
  env_var = {
    "SERVER" : "AWS",
    "APPLICATION" : "JOB"
    "EMAIL_SENDER" : var.email_sender,
    "TABLE_MEMBER" : module.dynamodb_tables.member_table_name,
    "TABLE_USER" : module.dynamodb_tables.user_table_name,
    "TABLE_CHURCH" : module.dynamodb_tables.church_table_name,
    "TABLE_PARTICIPANT" : module.dynamodb_tables.participant_table_name,
    "STORAGE" : module.s3_bucket.bucket_name,
    "QUEUE_URL" : module.sqs_queue.queue_url,
  }
}

module "job_schedule" {
  source          = "../../application/job_schedule"
  lambda_job_arn  = module.job.lambda_job_arn
  lambda_job_name = module.job.lambda_job_name
}

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