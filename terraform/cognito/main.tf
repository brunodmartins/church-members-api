resource "aws_cognito_user_pool" "user_pool" {
  name = "church-members-user-pool"
  admin_create_user_config {
    allow_admin_create_user_only = true
  }
}

output "user_pool_arn" {
  value = aws_cognito_user_pool.user_pool.arn
}