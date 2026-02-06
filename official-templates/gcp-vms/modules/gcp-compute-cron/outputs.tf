output "instance_id" {
  description = "The ID of the cron VM"
  value       = google_compute_instance.main.id
}
