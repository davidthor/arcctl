output "arn" {
  description = "Target group ARN"
  value       = aws_lb_target_group.this.arn
}

output "name" {
  description = "Target group name"
  value       = aws_lb_target_group.this.name
}
