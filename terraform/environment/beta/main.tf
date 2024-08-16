locals {
  ecr_image_tag       = "latest"
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
  source = "../../storage/dynamodb"
  member_table_name = var.member_table_name
  user_table_name = var.user_table_name
  church_table_name = var.church_table_name
}

module "iam_roles" {
  source = "../../security/iam"
  bucket_arn = module.s3_bucket.bucket_arn
  dynamodb_tables = module.dynamodb_tables.tables_arn
  role_name = var.role_name
}

module "s3_bucket" {
  source = "../../storage/s3"
  bucket_name = var.bucket_name
}


module "api" {
  source = "../../application/api"
  lambda_api_name = var.lambda_api_name
  lambda_role_arn = module.iam_roles.lambda_role_arn
  image_uri = "${data.aws_ecr_repository.repository.repository_url}@${data.aws_ecr_image.lambda_image.id}"
  env_var = {
    "SERVER" : "AWS",
    "APPLICATION" : "API"
    "EMAIL_SENDER": var.email_sender,
    "TABLE_MEMBER" : module.dynamodb_tables.member_table_name,
    "TABLE_USER" : module.dynamodb_tables.user_table_name,
    "TABLE_CHURCH": module.dynamodb_tables.church_table_name,
    "TOKEN_SECRET" : var.security_token_secret,
    "TOKEN_EXPIRATION" : var.security_token_expiration,
    "STORAGE": module.s3_bucket.bucket_name,
    "API_ENDPOINT": module.api.gateway_url
  }
  gateway_name    = var.gateway_name
}

output "gateway_d" {
  value = module.api.gateway_id
}