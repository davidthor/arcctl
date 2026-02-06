output "job_id" {
  description = "App Platform app ID"
  value       = digitalocean_app.job.id
}

output "status" {
  description = "Job status"
  value       = digitalocean_app.job.active_deployment_id != "" ? "complete" : "pending"
}
