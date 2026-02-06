output "host" {
  description = "Database host (private)"
  value       = digitalocean_database_cluster.db.private_host
}

output "port" {
  description = "Database port"
  value       = digitalocean_database_cluster.db.port
}

output "database" {
  description = "Default database name"
  value       = digitalocean_database_cluster.db.database
}

output "username" {
  description = "Database admin username"
  value       = digitalocean_database_cluster.db.user
}

output "password" {
  description = "Database admin password"
  value       = digitalocean_database_cluster.db.password
  sensitive   = true
}

output "connection_url" {
  description = "Full connection URL (using private host)"
  value       = local.connection_url
  sensitive   = true
}

output "cluster_id" {
  description = "Database cluster ID"
  value       = digitalocean_database_cluster.db.id
}
