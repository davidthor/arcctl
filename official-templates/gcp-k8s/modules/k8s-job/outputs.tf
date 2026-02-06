output "job_id" {
  description = "The name of the Kubernetes job"
  value       = kubernetes_job_v1.main.metadata[0].name
}

output "status" {
  description = "Job status"
  value       = "created"
}
