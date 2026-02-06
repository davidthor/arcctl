output "cronjob_id" {
  description = "CronJob ID"
  value       = kubernetes_cron_job_v1.this.metadata[0].uid
}

output "cronjob_name" {
  description = "CronJob name"
  value       = kubernetes_cron_job_v1.this.metadata[0].name
}
