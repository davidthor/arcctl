output "url" {
  description = "Route URL"
  value       = "https://${local.hostname}"
}

output "host" {
  description = "Route hostname"
  value       = local.hostname
}

output "port" {
  description = "Route port"
  value       = 443
}

output "fqdn" {
  description = "Fully qualified domain name"
  value       = digitalocean_record.route.fqdn
}
