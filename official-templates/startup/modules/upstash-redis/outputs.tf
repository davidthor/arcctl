output "host" {
  description = "Redis host"
  value       = upstash_redis_database.this.endpoint
}

output "port" {
  description = "Redis port"
  value       = upstash_redis_database.this.port
}

output "username" {
  description = "Redis username"
  value       = "default"
}

output "password" {
  description = "Redis password"
  value       = upstash_redis_database.this.password
  sensitive   = true
}

output "url" {
  description = "Full Redis connection URL"
  value       = "rediss://default:${upstash_redis_database.this.password}@${upstash_redis_database.this.endpoint}:${upstash_redis_database.this.port}"
  sensitive   = true
}
