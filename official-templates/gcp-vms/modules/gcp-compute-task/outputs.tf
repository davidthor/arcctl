output "instance_id" {
  description = "The ID of the task VM"
  value       = google_compute_instance.main.id
}

output "status" {
  description = "Task status"
  value       = "created"
}
