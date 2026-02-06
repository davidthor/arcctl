output "id" {
  description = "URL map ID for attaching backends"
  value       = google_compute_url_map.main.id
}

output "ip_address" {
  description = "Global static IP address"
  value       = google_compute_global_address.main.address
}
