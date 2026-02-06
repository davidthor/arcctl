output "cluster_name" {
  description = "The name of the GKE cluster"
  value       = google_container_cluster.main.name
}

output "endpoint" {
  description = "The endpoint of the GKE cluster"
  value       = google_container_cluster.main.endpoint
}

output "kubeconfig" {
  description = "Kubeconfig for connecting to the cluster"
  value = {
    host                   = "https://${google_container_cluster.main.endpoint}"
    cluster_ca_certificate = google_container_cluster.main.master_auth[0].cluster_ca_certificate
    token                  = data.google_client_config.current.access_token
  }
  sensitive = true
}

output "ingress_ip" {
  description = "The cluster endpoint IP for DNS records"
  value       = google_container_cluster.main.endpoint
}

data "google_client_config" "current" {}
