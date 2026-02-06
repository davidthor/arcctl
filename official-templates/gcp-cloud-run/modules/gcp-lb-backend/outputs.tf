output "backend_service_id" {
  description = "The ID of the backend service"
  value       = google_compute_backend_service.main.id
}

output "neg_id" {
  description = "The ID of the network endpoint group"
  value       = google_compute_region_network_endpoint_group.main.id
}
