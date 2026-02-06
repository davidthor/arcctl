output "username" {
  description = "Database username"
  value       = google_sql_user.user.name
}

output "password" {
  description = "Database password"
  value       = random_password.user.result
  sensitive   = true
}

output "connection_url" {
  description = "Full connection URL for this user"
  value       = "${local.scheme}://${google_sql_user.user.name}:${random_password.user.result}@${local.host}:${local.port}/${local.database_name}"
  sensitive   = true
}
