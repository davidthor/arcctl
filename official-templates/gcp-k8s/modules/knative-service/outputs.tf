output "service_id" {
  description = "The name of the Knative service"
  value       = var.name
}

output "url" {
  description = "The URL of the Knative service"
  value       = "http://${var.name}.${var.namespace}.svc.cluster.local"
}
