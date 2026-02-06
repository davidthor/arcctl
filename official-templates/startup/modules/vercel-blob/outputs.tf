output "endpoint" {
  description = "Blob store endpoint URL"
  value       = try(local.blob_data.url, "https://blob.vercel-storage.com")
}

output "bucket" {
  description = "Blob store name/bucket identifier"
  value       = random_id.store.hex
}

output "access_key_id" {
  description = "Access key for the blob store (store token)"
  value       = try(local.blob_data.clientToken, "")
  sensitive   = true
}

output "secret_access_key" {
  description = "Secret key for the blob store (store read-write token)"
  value       = try(local.blob_data.readWriteToken, "")
  sensitive   = true
}
