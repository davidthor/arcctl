output "id" {
  description = "Security group ID"
  value       = aws_security_group.this.id
}

output "subnet_id" {
  description = "First available private subnet ID"
  value       = local.subnet_ids[0]
}

output "subnet_ids" {
  description = "All available private subnet IDs"
  value       = local.subnet_ids
}
