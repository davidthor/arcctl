output "host" {
  description = "Internal hostname for the service"
  value       = local.host
}

output "port" {
  description = "Service port"
  value       = local.port
}

output "url" {
  description = "Full service URL"
  value       = local.service_url
}
