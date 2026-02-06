output "host" {
  description = "SMTP host"
  value       = var.host
}

output "port" {
  description = "SMTP port"
  value       = var.port
}

output "username" {
  description = "SMTP username"
  value       = var.username
}

output "password" {
  description = "SMTP password"
  value       = var.password
  sensitive   = true
}
