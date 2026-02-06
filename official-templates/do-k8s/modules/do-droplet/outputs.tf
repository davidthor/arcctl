output "droplet_id" {
  description = "Droplet ID"
  value       = digitalocean_droplet.droplet.id
}

output "ipv4_address" {
  description = "Droplet public IPv4 address"
  value       = digitalocean_droplet.droplet.ipv4_address
}
