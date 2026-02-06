output "bucket_name" {
  description = "S3 bucket name"
  value       = aws_s3_bucket.this.id
}

output "bucket_arn" {
  description = "S3 bucket ARN"
  value       = aws_s3_bucket.this.arn
}

output "access_key_id" {
  description = "IAM access key ID for bucket access"
  value       = aws_iam_access_key.this.id
  sensitive   = true
}

output "secret_access_key" {
  description = "IAM secret access key for bucket access"
  value       = aws_iam_access_key.this.secret
  sensitive   = true
}
