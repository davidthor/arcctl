output "service_id" {
  description = "The ID of the Cloud Run service"
  value       = google_cloud_run_v2_service.main.id
}

output "url" {
  description = "The URL of the Cloud Run service"
  value       = google_cloud_run_v2_service.main.uri
}
