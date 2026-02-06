output "image_uri" {
  description = "Full ECR image URI with tag"
  value       = local.image_uri
}

output "repository_url" {
  description = "ECR repository URL"
  value       = aws_ecr_repository.this.repository_url
}

output "repository_arn" {
  description = "ECR repository ARN"
  value       = aws_ecr_repository.this.arn
}
