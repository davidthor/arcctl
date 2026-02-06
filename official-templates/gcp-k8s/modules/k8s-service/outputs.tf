output "cluster_ip" {
  description = "ClusterIP address of the service"
  value       = kubernetes_service_v1.main.spec[0].cluster_ip
}

output "port" {
  description = "Service port"
  value       = var.port
}
