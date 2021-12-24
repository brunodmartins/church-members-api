resource "aws_dynamodb_table" "member_table" {
  name = "member"
  billing_mode = "PAY_PER_REQUEST"
  hash_key = "id"

  attribute {
    name = "id"
    type = "S"
  }
}

resource "aws_dynamodb_table" "member_v2" {
  name = "member_v2"
  billing_mode = "PAY_PER_REQUEST"
  hash_key = "church_id"
  range_key = "id"

  attribute {
    name = "church_id"
    type = "S"
  }

  attribute {
    name = "id"
    type = "S"
  }

  attribute {
    name = "birthDateShort"
    type = "S"
  }

  attribute {
    name = "maritalStatus"
    type = "S"
  }

  attribute {
    name = "name"
    type = "S"
  }

  global_secondary_index {
    name = "birthDateIndex"
    hash_key = "church_id"
    range_key = "birthDateShort"
    projection_type = "INCLUDE"
    non_key_attributes = [
      "id",
      "church_id",
      "active",
      "birthDate",
      "firstName",
      "lastName",
      "name",
      "gender",
      "marriageDate"]
  }

  global_secondary_index {
    name = "maritalStatusIndex"
    hash_key = "church_id"
    range_key = "maritalStatus"
    projection_type = "INCLUDE"
    non_key_attributes = [
      "id",
      "church_id",
      "active",
      "birthDate",
      "firstName",
      "lastName",
      "name",
      "spousesName",
      "gender",
      "marriageDate",
      "marriageDateShort"]
  }

  global_secondary_index {
    name = "nameIndex"
    hash_key = "church_id"
    range_key = "name"
    projection_type = "INCLUDE"
    non_key_attributes = [
      "id",
      "church_id",
      "active",
      "birthDate",
      "firstName",
      "lastName",
      "name",
      "gender",
      "marriageDate"]
  }

}

resource "aws_dynamodb_table" "member_history_table" {
  name = "member_history"
  billing_mode = "PAY_PER_REQUEST"
  hash_key = "id"

  attribute {
    name = "id"
    type = "S"
  }
}

resource "aws_dynamodb_table" "user_table" {
  name = "user"
  billing_mode = "PAY_PER_REQUEST"
  hash_key = "id"

  attribute {
    name = "id"
    type = "S"
  }
}

resource "aws_dynamodb_table" "users_v2" {
  name = "user_v2"
  billing_mode = "PAY_PER_REQUEST"
  hash_key = "church_id"
  range_key = "username"

  attribute {
    name = "church_id"
    type = "S"
  }

  attribute {
    name = "username"
    type = "S"
  }
}


resource "aws_dynamodb_table" "church_table" {
  name = "church"
  billing_mode = "PAY_PER_REQUEST"
  hash_key = "id"

  attribute {
    name = "id"
    type = "S"
  }
}

output "member_table_name" {
  value = aws_dynamodb_table.member_v2.name
}

output "member_history_table_name" {
  value = aws_dynamodb_table.member_history_table.name
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
    aws_dynamodb_table.member_history_table.arn,
    aws_dynamodb_table.users_v2.arn,
    aws_dynamodb_table.church_table.arn
  ]
}

