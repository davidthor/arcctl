output "instance_id" {
  description = "EC2 instance ID"
  value       = aws_instance.this.id
}

output "private_ip" {
  description = "Private IP address"
  value       = aws_instance.this.private_ip
}

output "port" {
  description = "Application port"
  value       = local.app_port
}

output "public_ip" {
  description = "Public IP address"
  value       = aws_instance.this.public_ip
}
