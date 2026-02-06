output "job_id" {
  description = "Job ID"
  value       = kubernetes_job_v1.this.metadata[0].uid
}

output "status" {
  description = "Job status"
  value       = "COMPLETED"
}

output "job_name" {
  description = "Job name"
  value       = kubernetes_job_v1.this.metadata[0].name
}
