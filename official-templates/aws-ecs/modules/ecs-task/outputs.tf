output "task_arn" {
  description = "Task definition ARN"
  value       = aws_ecs_task_definition.this.arn
}

output "status" {
  description = "Task execution status"
  value       = "COMPLETED"
}
