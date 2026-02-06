output "gateway_name" {
  description = "Name of the created Gateway"
  value       = kubernetes_manifest.gateway.manifest.metadata.name
}
