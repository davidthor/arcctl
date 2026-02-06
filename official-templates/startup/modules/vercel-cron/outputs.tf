output "cron_id" {
  description = "Cron job identifier"
  value       = "${var.project_id}-cron-${var.name}"
}
