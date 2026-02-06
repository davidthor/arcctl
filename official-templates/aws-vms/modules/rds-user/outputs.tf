output "username" {
  description = "Database username"
  value       = local.username
}

output "password" {
  description = "Database password"
  value       = local.password
  sensitive   = true
}

output "connection_url" {
  description = "Database connection URL with new user credentials"
  value       = local.connection_url
  sensitive   = true
}
