output "job_id" {
  description = "The ID of the Cloud Scheduler job"
  value       = google_cloud_scheduler_job.main.id
}
