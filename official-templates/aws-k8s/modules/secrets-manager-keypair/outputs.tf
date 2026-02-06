output "private_key" {
  description = "Private key in PEM format"
  value       = tls_private_key.this.private_key_pem
  sensitive   = true
}

output "public_key" {
  description = "Public key in PEM format"
  value       = tls_private_key.this.public_key_pem
  sensitive   = true
}

output "private_key_base64" {
  description = "Private key in base64-encoded PEM format"
  value       = base64encode(tls_private_key.this.private_key_pem)
  sensitive   = true
}

output "public_key_base64" {
  description = "Public key in base64-encoded PEM format"
  value       = base64encode(tls_private_key.this.public_key_pem)
  sensitive   = true
}
