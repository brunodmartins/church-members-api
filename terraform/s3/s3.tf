
resource "aws_s3_bucket" "bucket" {
  bucket = "church-members-data"
}

resource "aws_s3_bucket_server_side_encryption_configuration" "sse_config" {
  bucket = aws_s3_bucket.bucket.id
  rule {
      apply_server_side_encryption_by_default {
        sse_algorithm     = "aws:kms"
      }
  }
}

resource "aws_s3_bucket_lifecycle_configuration" "lifecycle" {
  bucket = aws_s3_bucket.bucket.id
  
  rule {
    id = "expire-failed-uploads"
    status = "Enabled"
    abort_incomplete_multipart_upload {
        days_after_initiation = 1
    }
  }
}

resource "aws_s3_bucket_versioning" "versioning" {
  bucket = aws_s3_bucket.bucket.id
  versioning_configuration {
    status = "Enabled"
  }
}


resource "aws_s3_bucket_acl" "acl" {
  bucket = aws_s3_bucket.bucket.id
  acl    = "private"
}

resource "aws_s3_bucket_public_access_block" "block_public_access" {
  bucket = aws_s3_bucket.bucket.id

  block_public_acls   = true
  block_public_policy = true
  ignore_public_acls = true
  restrict_public_buckets = true
}


output "bucket_name" {
    value = "church-members-data"
}

output "bucket_arn" {
    value = aws_s3_bucket.bucket.arn
}