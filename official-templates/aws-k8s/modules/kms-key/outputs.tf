output "key_id" {
  description = "KMS key ID"
  value       = aws_kms_key.this.key_id
}

output "key_arn" {
  description = "KMS key ARN"
  value       = aws_kms_key.this.arn
}

output "key_material" {
  description = "Raw symmetric key material (hex-encoded)"
  value       = random_bytes.key_material.hex
  sensitive   = true
}

output "key_material_base64" {
  description = "Symmetric key material (base64-encoded)"
  value       = random_bytes.key_material.base64
  sensitive   = true
}
