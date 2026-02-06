# Asymmetric key outputs
output "privateKey" {
  description = "PEM-encoded private key"
  value       = var.key_type != "symmetric" ? tls_private_key.asymmetric[0].private_key_pem : ""
  sensitive   = true
}

output "publicKey" {
  description = "PEM-encoded public key"
  value       = var.key_type != "symmetric" ? tls_private_key.asymmetric[0].public_key_pem : ""
}

output "privateKeyBase64" {
  description = "Base64-encoded private key"
  value       = var.key_type != "symmetric" ? base64encode(tls_private_key.asymmetric[0].private_key_pem) : ""
  sensitive   = true
}

output "publicKeyBase64" {
  description = "Base64-encoded public key"
  value       = var.key_type != "symmetric" ? base64encode(tls_private_key.asymmetric[0].public_key_pem) : ""
}

# Symmetric key outputs
output "key" {
  description = "Raw symmetric key (hex)"
  value       = var.key_type == "symmetric" ? random_bytes.symmetric[0].hex : ""
  sensitive   = true
}

output "keyBase64" {
  description = "Base64-encoded symmetric key"
  value       = var.key_type == "symmetric" ? random_bytes.symmetric[0].base64 : ""
  sensitive   = true
}
