output "fqdn" {
  description = "Fully qualified domain name"
  value       = digitalocean_record.record.fqdn
}
