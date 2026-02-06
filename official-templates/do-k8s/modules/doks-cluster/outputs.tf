output "kubeconfig" {
  description = "Kubernetes cluster kubeconfig"
  value       = digitalocean_kubernetes_cluster.cluster.kube_config[0].raw_config
  sensitive   = true
}

output "cluster_id" {
  description = "Kubernetes cluster ID"
  value       = digitalocean_kubernetes_cluster.cluster.id
}

output "endpoint" {
  description = "Kubernetes API server endpoint"
  value       = digitalocean_kubernetes_cluster.cluster.endpoint
}

output "load_balancer_ip" {
  description = "Load balancer external IP from ingress controller"
  value       = try(data.kubernetes_service_v1.nginx_lb.status[0].load_balancer[0].ingress[0].ip, "")
}
