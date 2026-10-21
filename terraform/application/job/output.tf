output "lambda_job_arn" {
  value = aws_lambda_function.lambda_job.arn
}

output "lambda_job_name" {
  value = aws_lambda_function.lambda_job.function_name
}
