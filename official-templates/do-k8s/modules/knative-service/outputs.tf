output "service_id" {
  description = "Knative service name"
  value       = kubernetes_manifest.knative_service.manifest.metadata.name
}

output "url" {
  description = "Knative service URL"
  value       = try(kubernetes_manifest.knative_service.object.status.url, "http://${local.name}.${var.namespace}.svc.cluster.local")
}
