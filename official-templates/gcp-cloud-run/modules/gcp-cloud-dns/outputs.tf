output "fqdn" {
  description = "Fully qualified domain name"
  value       = google_dns_record_set.main.name
}
