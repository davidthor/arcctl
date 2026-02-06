variable "name" {
  description = "Name of the Gateway"
  type        = string
}

variable "namespace" {
  description = "Kubernetes namespace"
  type        = string
}

variable "gateway_class" {
  description = "GatewayClass name (e.g., gke-l7-global-external-managed)"
  type        = string
}

variable "kubeconfig" {
  description = "Kubeconfig for the Kubernetes cluster"
  type = object({
    host                   = string
    cluster_ca_certificate = string
    token                  = string
  })
  sensitive = true
}

variable "tls" {
  description = "TLS configuration"
  type = object({
    enabled = bool
    issuer  = optional(string)
  })
  default = {
    enabled = false
  }
}
