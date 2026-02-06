output "cluster_ip" {
  description = "Service ClusterIP"
  value       = kubernetes_service_v1.this.spec[0].cluster_ip
}

output "port" {
  description = "Service port"
  value       = kubernetes_service_v1.this.spec[0].port[0].port
}

output "name" {
  description = "Service name"
  value       = kubernetes_service_v1.this.metadata[0].name
}
