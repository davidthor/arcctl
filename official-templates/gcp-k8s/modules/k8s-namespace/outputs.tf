output "name" {
  description = "The name of the created namespace"
  value       = kubernetes_namespace_v1.main.metadata[0].name
}
