output "private_ip" {
  description = "Private IP address of the Cloud SQL instance"
  value       = google_sql_database_instance.main.private_ip_address
}

output "port" {
  description = "Database port"
  value       = local.port
}

output "database" {
  description = "Name of the database"
  value       = google_sql_database.main.name
}

output "username" {
  description = "Admin username"
  value       = google_sql_user.admin.name
}

output "password" {
  description = "Admin password"
  value       = random_password.admin.result
  sensitive   = true
}

output "connection_url" {
  description = "Full connection URL"
  value       = "${local.scheme}://${google_sql_user.admin.name}:${random_password.admin.result}@${google_sql_database_instance.main.private_ip_address}:${local.port}/${google_sql_database.main.name}"
  sensitive   = true
}

output "instance_name" {
  description = "Cloud SQL instance name"
  value       = google_sql_database_instance.main.name
}
