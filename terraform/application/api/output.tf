output "gateway_id" {
  value = aws_api_gateway_rest_api.api_gateway.id
}

output "gateway_url" {
  value = aws_api_gateway_stage.api_stage.invoke_url
}