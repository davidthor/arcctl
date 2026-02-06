# Asymmetric key outputs (RSA / ECDSA)
output "privateKey" {
  description = "Private key in PEM format"
  value       = local.private_key_pem
  sensitive   = true
}

output "publicKey" {
  description = "Public key in PEM format"
  value       = local.public_key_pem
}

output "privateKeyBase64" {
  description = "Private key encoded as base64"
  value       = local.private_key_pem != "" ? base64encode(local.private_key_pem) : ""
  sensitive   = true
}

output "publicKeyBase64" {
  description = "Public key encoded as base64"
  value       = local.public_key_pem != "" ? base64encode(local.public_key_pem) : ""
}

# Symmetric key outputs
output "key" {
  description = "Symmetric key as hex string"
  value       = local.symmetric_key
  sensitive   = true
}

output "keyBase64" {
  description = "Symmetric key encoded as base64"
  value       = local.is_sym ? try(random_bytes.symmetric[0].base64, "") : ""
  sensitive   = true
}
