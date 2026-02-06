output "stage_id" {
  description = "Stage ID"
  value       = aws_apigatewayv2_stage.this.id
}

output "invoke_url" {
  description = "Stage invoke URL"
  value       = aws_apigatewayv2_stage.this.invoke_url
}
