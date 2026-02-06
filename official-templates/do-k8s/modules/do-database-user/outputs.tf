output "username" {
  description = "Database username"
  value       = digitalocean_database_user.user.name
}

output "password" {
  description = "Database password"
  value       = digitalocean_database_user.user.password
  sensitive   = true
}

output "url" {
  description = "Connection URL"
  value       = local.connection_url
  sensitive   = true
}
