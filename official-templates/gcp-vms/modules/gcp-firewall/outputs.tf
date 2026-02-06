output "id" {
  description = "The ID of the main firewall rule"
  value       = google_compute_firewall.allow_http.id
}
