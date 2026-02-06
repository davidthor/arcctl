output "route_name" {
  description = "HTTPRoute name"
  value       = kubernetes_manifest.httproute.manifest.metadata.name
}
