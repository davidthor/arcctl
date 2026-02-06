output "backend_service_id" {
  description = "The ID of the backend service"
  value       = google_compute_backend_service.main.id
}

output "instance_group_id" {
  description = "The ID of the instance group"
  value       = google_compute_instance_group.main.id
}
