output "cluster_ip" {
  description = "ClusterIP address"
  value       = kubernetes_service_v1.service.spec[0].cluster_ip
}

output "port" {
  description = "Service port"
  value       = kubernetes_service_v1.service.spec[0].port[0].port
}
