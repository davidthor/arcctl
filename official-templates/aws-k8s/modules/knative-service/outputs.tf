output "service_id" {
  description = "Knative service ID"
  value       = kubernetes_manifest.knative_service.manifest.metadata.name
}

output "url" {
  description = "Knative service URL"
  value       = "http://${local.name}.${var.namespace}.svc.cluster.local"
}
