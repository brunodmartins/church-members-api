
module "dynamodb_tables" {
  source = "../../storage/dynamodb"
  member_table_name = var.member_table_name
  user_table_name = var.user_table_name
  church_table_name = var.church_table_name
  participant_table_name = var.participant_table_name
}

module "s3_bucket" {
  source = "../../storage/s3"
  bucket_name = var.bucket_name
}