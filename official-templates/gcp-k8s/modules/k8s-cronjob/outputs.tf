output "cronjob_id" {
  description = "The name of the CronJob"
  value       = kubernetes_cron_job_v1.main.metadata[0].name
}
