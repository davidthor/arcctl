output "deployment_id" {
  description = "Deployment UID"
  value       = kubernetes_deployment_v1.deployment.metadata[0].uid
}

output "deployment_name" {
  description = "Deployment name"
  value       = kubernetes_deployment_v1.deployment.metadata[0].name
}
