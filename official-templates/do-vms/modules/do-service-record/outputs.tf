output "host" {
  description = "Service hostname"
  value       = digitalocean_record.service.fqdn
}

output "port" {
  description = "Service port"
  value       = local.service_port
}

output "url" {
  description = "Service URL"
  value       = "http://${digitalocean_record.service.fqdn}:${local.service_port}"
}
