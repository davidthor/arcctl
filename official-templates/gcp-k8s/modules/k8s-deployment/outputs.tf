output "deployment_id" {
  description = "The name of the Kubernetes deployment"
  value       = kubernetes_deployment_v1.main.metadata[0].name
}
