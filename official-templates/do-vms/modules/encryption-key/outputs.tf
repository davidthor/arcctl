# Asymmetric key outputs
output "privateKey" {
  description = "Private key in PEM format"
  value       = var.key_type == "rsa" ? try(tls_private_key.rsa[0].private_key_pem, "") : try(tls_private_key.ecdsa[0].private_key_pem, "")
  sensitive   = true
}

output "publicKey" {
  description = "Public key in PEM format"
  value       = var.key_type == "rsa" ? try(tls_private_key.rsa[0].public_key_pem, "") : try(tls_private_key.ecdsa[0].public_key_pem, "")
}

output "privateKeyBase64" {
  description = "Private key base64 encoded"
  value       = var.key_type == "rsa" ? try(base64encode(tls_private_key.rsa[0].private_key_pem), "") : try(base64encode(tls_private_key.ecdsa[0].private_key_pem), "")
  sensitive   = true
}

output "publicKeyBase64" {
  description = "Public key base64 encoded"
  value       = var.key_type == "rsa" ? try(base64encode(tls_private_key.rsa[0].public_key_pem), "") : try(base64encode(tls_private_key.ecdsa[0].public_key_pem), "")
}

# Symmetric key outputs
output "key" {
  description = "Symmetric key in hex format"
  value       = try(random_bytes.symmetric[0].hex, "")
  sensitive   = true
}

output "keyBase64" {
  description = "Symmetric key base64 encoded"
  value       = try(random_bytes.symmetric[0].base64, "")
  sensitive   = true
}
