output "bucket_name" {
  value = var.bucket_name
}

output "bucket_arn" {
  value = aws_s3_bucket.bucket.arn
}