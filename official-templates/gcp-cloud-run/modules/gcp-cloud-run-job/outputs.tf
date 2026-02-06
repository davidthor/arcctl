output "job_id" {
  description = "The ID of the Cloud Run job"
  value       = google_cloud_run_v2_job.main.id
}

output "status" {
  description = "Status of the Cloud Run job"
  value       = "created"
}
