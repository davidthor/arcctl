terraform {
  required_providers {
    digitalocean = {
      source  = "digitalocean/digitalocean"
      version = "~> 2.0"
    }
  }
}

provider "digitalocean" {
  token = var.token
}

resource "digitalocean_kubernetes_cluster" "cluster" {
  name    = var.name
  region  = var.region
  version = data.digitalocean_kubernetes_versions.current.latest_version

  node_pool {
    name       = "${var.name}-default"
    size       = var.node_pool.size
    auto_scale = var.node_pool.auto_scale
    min_nodes  = var.node_pool.min_nodes
    max_nodes  = var.node_pool.max_nodes

    tags = ["arcctl", var.name]
  }

  tags = ["arcctl", "managed-by:arcctl"]
}

data "digitalocean_kubernetes_versions" "current" {
  version_prefix = "1."
}

# Install NGINX Ingress Controller for Gateway API support
resource "helm_release" "nginx_ingress" {
  name             = "ingress-nginx"
  repository       = "https://kubernetes.github.io/ingress-nginx"
  chart            = "ingress-nginx"
  namespace        = "ingress-nginx"
  create_namespace = true
  version          = "4.9.0"

  set {
    name  = "controller.service.type"
    value = "LoadBalancer"
  }

  set {
    name  = "controller.service.annotations.service\\.beta\\.kubernetes\\.io/do-loadbalancer-name"
    value = "${var.name}-lb"
  }

  depends_on = [digitalocean_kubernetes_cluster.cluster]
}

# Install cert-manager for TLS
resource "helm_release" "cert_manager" {
  name             = "cert-manager"
  repository       = "https://charts.jetstack.io"
  chart            = "cert-manager"
  namespace        = "cert-manager"
  create_namespace = true
  version          = "1.14.0"

  set {
    name  = "installCRDs"
    value = "true"
  }

  depends_on = [digitalocean_kubernetes_cluster.cluster]
}

# Wait for the load balancer to get an IP
data "kubernetes_service_v1" "nginx_lb" {
  metadata {
    name      = "ingress-nginx-controller"
    namespace = "ingress-nginx"
  }

  depends_on = [helm_release.nginx_ingress]
}
