output "deployment_id" {
  description = "Deployment ID"
  value       = kubernetes_deployment_v1.this.metadata[0].uid
}

output "deployment_name" {
  description = "Deployment name"
  value       = kubernetes_deployment_v1.this.metadata[0].name
}
