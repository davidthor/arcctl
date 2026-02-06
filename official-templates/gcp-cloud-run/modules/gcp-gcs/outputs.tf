output "bucket_name" {
  description = "Name of the GCS bucket"
  value       = google_storage_bucket.main.name
}

output "hmac_access_key" {
  description = "HMAC access key ID for S3-compatible access"
  value       = google_storage_hmac_key.main.access_id
}

output "hmac_secret_key" {
  description = "HMAC secret key for S3-compatible access"
  value       = google_storage_hmac_key.main.secret
  sensitive   = true
}
