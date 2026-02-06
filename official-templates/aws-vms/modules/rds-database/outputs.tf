output "endpoint" {
  description = "Database endpoint hostname"
  value       = local.endpoint
}

output "port" {
  description = "Database port"
  value       = local.port
}

output "database_name" {
  description = "Database name"
  value       = local.database_name
}

output "username" {
  description = "Database master username"
  value       = local.username
  sensitive   = true
}

output "password" {
  description = "Database master password"
  value       = local.password
  sensitive   = true
}

output "connection_url" {
  description = "Full database connection URL"
  value       = local.connection_url
  sensitive   = true
}
