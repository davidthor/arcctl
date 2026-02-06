output "droplet_id" {
  description = "Droplet ID"
  value       = digitalocean_droplet.task.id
}

output "status" {
  description = "Task status"
  value       = digitalocean_droplet.task.status
}
