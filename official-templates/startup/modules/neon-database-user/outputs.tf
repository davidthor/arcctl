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
  value       = "postgresql://${neon_role.this.name}:${neon_role.this.password}@${local.host}:5432/${var.database}?sslmode=require"
  sensitive   = true
}
