output "dns_name" {
  description = "Service discovery DNS name"
  value       = "${var.name}.${var.namespace}.local"
}

output "port" {
  description = "Service port"
  value       = var.port
}

output "service_arn" {
  description = "Service discovery service ARN"
  value       = aws_service_discovery_service.this.arn
}

output "namespace_id" {
  description = "Service discovery namespace ID"
  value       = aws_service_discovery_private_dns_namespace.this.id
}
