output "cronjob_id" {
  description = "CronJob UID"
  value       = kubernetes_cron_job_v1.cronjob.metadata[0].uid
}
