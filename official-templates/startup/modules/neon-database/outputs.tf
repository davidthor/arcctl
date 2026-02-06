output "host" {
  description = "Database host"
  value       = neon_endpoint.this.host
}

output "port" {
  description = "Database port"
  value       = 5432
}

output "database" {
  description = "Database name"
  value       = neon_database.this.name
}

output "username" {
  description = "Database username"
  value       = neon_role.this.name
}

output "password" {
  description = "Database password"
  value       = neon_role.this.password
  sensitive   = true
}

output "url" {
  description = "Full PostgreSQL connection URL"
  value       = "postgresql://${neon_role.this.name}:${neon_role.this.password}@${neon_endpoint.this.host}:5432/${neon_database.this.name}?sslmode=require"
  sensitive   = true
}
