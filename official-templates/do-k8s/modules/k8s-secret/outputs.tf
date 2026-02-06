output "secret_name" {
  description = "Name of the created secret"
  value       = kubernetes_secret_v1.secret.metadata[0].name
}

output "id" {
  description = "Secret UID"
  value       = kubernetes_secret_v1.secret.metadata[0].uid
}
