output "url" {
  description = "Route URL"
  value       = "https://${var.domain}"
}

output "hosts" {
  description = "Route hosts"
  value       = [var.domain]
}

output "host" {
  description = "Route host"
  value       = var.domain
}

output "port" {
  description = "Route port"
  value       = 443
}
