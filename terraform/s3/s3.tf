
resource "aws_s3_bucket" "bucket" {
  bucket = "church-members-data"
  acl    = "private"

  versioning {
    enabled = true
  }

  server_side_encryption_configuration {
    rule {
      apply_server_side_encryption_by_default {
        sse_algorithm     = "aws:kms"
      }
    }
  }

  lifecycle_rule {
    id = "expire-failed-uploads"
    enabled = true

    abort_incomplete_multipart_upload_days = 1
  }
}

output "bucket_name" {
    value = "church-members-data"
}

output "bucket_arn" {
    value = aws_s3_bucket.bucket.arn
}