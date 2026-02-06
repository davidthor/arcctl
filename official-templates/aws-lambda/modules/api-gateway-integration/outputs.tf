output "integration_id" {
  description = "Integration ID"
  value       = aws_apigatewayv2_integration.this.id
}

output "endpoint_host" {
  description = "Service endpoint host"
  value       = local.endpoint_host
}

output "endpoint_port" {
  description = "Service endpoint port"
  value       = var.port
}

output "endpoint_url" {
  description = "Service endpoint URL"
  value       = "${data.aws_apigatewayv2_api.this.api_endpoint}/${var.stage}/${var.name}"
}
