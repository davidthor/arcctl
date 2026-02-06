output "secret_id" {
  description = "Secret identifier"
  value       = random_id.secret_id.hex
}

output "data" {
  description = "Secret data"
  value       = var.data
  sensitive   = true
}
