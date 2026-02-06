output "bucket_name" {
  description = "Spaces bucket name"
  value       = digitalocean_spaces_bucket.bucket.name
}

output "endpoint" {
  description = "Spaces bucket endpoint"
  value       = digitalocean_spaces_bucket.bucket.bucket_domain_name
}

output "access_key" {
  description = "Spaces access key"
  value       = digitalocean_spaces_bucket.bucket.access_key_id
  sensitive   = true
}

output "secret_key" {
  description = "Spaces secret key"
  value       = digitalocean_spaces_bucket.bucket.secret_access_key
  sensitive   = true
}

output "urn" {
  description = "Spaces bucket URN"
  value       = digitalocean_spaces_bucket.bucket.urn
}
