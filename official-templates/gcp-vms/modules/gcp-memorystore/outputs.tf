output "host" {
  description = "Redis instance host"
  value       = google_redis_instance.main.host
}

output "port" {
  description = "Redis instance port"
  value       = google_redis_instance.main.port
}

output "connection_url" {
  description = "Redis connection URL"
  value       = "redis://${google_redis_instance.main.host}:${google_redis_instance.main.port}"
}
