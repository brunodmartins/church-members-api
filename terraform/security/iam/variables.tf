variable "dynamodb_tables" {
  type = list(string)
}

variable "bucket_arn" {
  type = string
}

variable "role_name" {
  type = string
}