output "instance_id" {
  description = "The ID of the Compute Engine instance"
  value       = google_compute_instance.main.id
}

output "internal_ip" {
  description = "Internal IP address"
  value       = google_compute_instance.main.network_interface[0].network_ip
}

output "external_ip" {
  description = "External IP address"
  value       = google_compute_instance.main.network_interface[0].access_config[0].nat_ip
}
