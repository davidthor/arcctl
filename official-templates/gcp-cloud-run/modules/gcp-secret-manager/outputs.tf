output "secret_id" {
  description = "The ID of the secret"
  value       = google_secret_manager_secret.main.secret_id
}
