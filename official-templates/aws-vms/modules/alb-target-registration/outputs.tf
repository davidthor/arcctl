output "id" {
  description = "Target group attachment ID"
  value       = aws_lb_target_group_attachment.this.id
}

output "dns_name" {
  description = "ALB DNS name for service discovery"
  value       = data.aws_lb.this.dns_name
}

output "port" {
  description = "Service port"
  value       = var.port
}
