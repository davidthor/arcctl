output "function_arn" {
  description = "Lambda function ARN"
  value       = aws_lambda_function.this.arn
}

output "status" {
  description = "Task execution status"
  value       = "COMPLETED"
}

output "function_name" {
  description = "Lambda function name"
  value       = aws_lambda_function.this.function_name
}
