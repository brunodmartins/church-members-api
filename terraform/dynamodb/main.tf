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

resource "aws_dynamodb_table" "member_v2" {
  name           = "member_v2"
  billing_mode   = "PROVISIONED"
  read_capacity  = 5
  write_capacity = 5
  hash_key       = "church_id"
  range_key      = "id"

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
    name = "marriageDateShort"
    type = "S"
  }

  attribute {
    name = "name"
    type = "S"
  }

  global_secondary_index {
    name               = "birthDateIndex"
    hash_key           = "church_id"
    range_key          = "birthDateShort"
    write_capacity     = 5
    read_capacity      = 5
    projection_type    = "INCLUDE"
    non_key_attributes = ["id","church_id","active","birthDate", "firstName", "lastName", "name", "gender", "marriageDate"]
  }

  global_secondary_index {
    name               = "marriageDateIndex"
    hash_key           = "church_id"
    range_key          = "marriageDateShort"
    write_capacity     = 5
    read_capacity      = 5
    projection_type    = "INCLUDE"
    non_key_attributes = ["id","church_id","active","birthDate", "firstName", "lastName", "name", "gender", "marriageDate"]
  }

  global_secondary_index {
    name               = "nameIndex"
    hash_key           = "church_id"
    range_key          = "name"
    write_capacity     = 5
    read_capacity      = 5
    projection_type    = "INCLUDE"
    non_key_attributes = ["id","church_id","active","birthDate", "firstName", "lastName", "name", "gender", "marriageDate"]
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

resource "aws_dynamodb_table" "user_table" {
  name           = "user"
  billing_mode   = "PROVISIONED"
  read_capacity  = 5
  write_capacity = 5
  hash_key       = "id"

  attribute {
    name = "id"
    type = "S"
  }
}

resource "aws_dynamodb_table" "users_v2" {
  name           = "user_v2"
  billing_mode   = "PROVISIONED"
  read_capacity  = 5
  write_capacity = 5
  hash_key       = "church_id"
  range_key      = "username"

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
  name           = "church"
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

output "user_table_name" {
  value = aws_dynamodb_table.user_table.name
}

output "church_table_name" {
  value = aws_dynamodb_table.church_table.name
}

output "tables_arn" {
  value = [
    aws_dynamodb_table.member_table.arn,
    aws_dynamodb_table.member_history_table.arn,
    aws_dynamodb_table.user_table.arn,
    aws_dynamodb_table.church_table.arn
  ]
}

