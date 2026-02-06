output "host" {
  description = "Internal service hostname"
  value       = local.internal_host
}

output "port" {
  description = "Internal service port"
  value       = local.service_port
}

output "url" {
  description = "Internal service URL"
  value       = "http://${local.internal_host}:${local.service_port}"
}
