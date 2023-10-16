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