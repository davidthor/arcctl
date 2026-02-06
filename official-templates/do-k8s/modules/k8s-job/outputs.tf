output "job_id" {
  description = "Job ID"
  value       = kubernetes_job_v1.job.metadata[0].uid
}

output "status" {
  description = "Job status"
  value       = "complete"
}
