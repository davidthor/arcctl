output "ipv4_address" {
  description = "OTel collector Droplet public IPv4 address"
  value       = digitalocean_droplet.otel.ipv4_address
}

output "ipv4_address_private" {
  description = "OTel collector Droplet private IPv4 address"
  value       = digitalocean_droplet.otel.ipv4_address_private
}

output "droplet_id" {
  description = "Droplet ID"
  value       = digitalocean_droplet.otel.id
}
