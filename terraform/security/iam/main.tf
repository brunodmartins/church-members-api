resource "aws_iam_policy" "church_members_api_policy" {
  name = "${var.role_name}-policy"
  description = "This policy allow church-members-api full execution"
  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "dynamodb:PutItem",
          "dynamodb:UpdateItem",
          "dynamodb:DeleteItem",
          "dynamodb:GetItem",
          "dynamodb:Query",
          "dynamodb:Scan",
        ]
        Resource = var.dynamodb_tables
      },
      {
        Effect = "Allow"
        Action = [
          "logs:CreateLogGroup",
          "logs:CreateLogStream",
          "logs:PutLogEvents",
        ]
        Resource = "*"
      },
      {
        Effect = "Allow"
        Action = [
          "sns:Publish",
        ]
        Resource = "*"
      },
      {
        Effect = "Allow"
        Action = [
          "SES:SendEmail",
          "SES:SendRawEmail"
        ]
        Resource = "*"
      },
      {
        Effect = "Allow"
        Action = [
          "s3:GetBucketLocation"
        ]
        Resource = var.bucket_arn
      },
      {
        Effect = "Allow"
        Action = [
          "s3:GetObject",
          "s3:PutObject"
        ]
        Resource = "${var.bucket_arn}/*"
      }

    ]
  })
}

resource "aws_iam_role" "church_members_api_role" {
  name = "${var.role_name}-role"
  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Principal = {
          Service = "lambda.amazonaws.com"
        }
      }
    ]
  })
}

resource "aws_iam_role_policy_attachment" "attach_policy" {
  role = aws_iam_role.church_members_api_role.name
  policy_arn = aws_iam_policy.church_members_api_policy.arn
}

output "lambda_role_arn" {
  value = aws_iam_role.church_members_api_role.arn
}