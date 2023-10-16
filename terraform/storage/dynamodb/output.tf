output "member_table_name" {
  value = aws_dynamodb_table.member_v2.name
}

output "user_table_name" {
  value = aws_dynamodb_table.users_v2.name
}

output "church_table_name" {
  value = aws_dynamodb_table.church_table.name
}

output "tables_arn" {
  value = [
    aws_dynamodb_table.member_v2.arn,
    "${aws_dynamodb_table.member_v2.arn}/index/nameIndex",
    "${aws_dynamodb_table.member_v2.arn}/index/maritalStatusIndex",
    "${aws_dynamodb_table.member_v2.arn}/index/birthDateIndex",
    aws_dynamodb_table.users_v2.arn,
    aws_dynamodb_table.church_table.arn
  ]
}