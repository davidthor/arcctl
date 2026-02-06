output "rule_arn" {
  description = "EventBridge rule ARN"
  value       = aws_cloudwatch_event_rule.this.arn
}

output "function_arn" {
  description = "Lambda function ARN"
  value       = aws_lambda_function.this.arn
}

output "id" {
  description = "Scheduled Lambda ID"
  value       = aws_cloudwatch_event_rule.this.id
}
