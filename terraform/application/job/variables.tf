variable "lambda_job_name" {
  type = string
}

variable "lambda_role_arn" {
  type = string
}

variable "image_uri" {
  type = string
}

variable "env_var" {
  type = map(string)
}

variable "reports_topic_name" {
  type    = string
  default = "weekly-report-topic"
}