resource "aws_dynamodb_table" "member_table" {
  name           = "member"
  billing_mode   = "PROVISIONED"
  read_capacity  = 5
  write_capacity = 5
  hash_key       = "id"

  attribute {
    name = "id"
    type = "S"
  }
}

resource "aws_dynamodb_table" "member_history_table" {
  name           = "member_history"
  billing_mode   = "PROVISIONED"
  read_capacity  = 5
  write_capacity = 5
  hash_key       = "id"

  attribute {
    name = "id"
    type = "S"
  }
}

output "member_table_name" {
  value = aws_dynamodb_table.member_table.name
}

output "member_history_table_name" {
  value = aws_dynamodb_table.member_history_table.name
}

output "tables_arn" {
  value = [aws_dynamodb_table.member_table.arn, aws_dynamodb_table.member_history_table.arn]
}

