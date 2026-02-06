variable "name" {
  description = "Gateway name"
  type        = string
}

variable "namespace" {
  description = "Kubernetes namespace"
  type        = string
}

variable "gateway_class" {
  description = "Gateway class name"
  type        = string
  default     = "nginx"
}

variable "kubeconfig" {
  description = "Kubernetes cluster kubeconfig content"
  type        = string
  sensitive   = true
}

variable "tls" {
  description = "TLS configuration"
  type = object({
    enabled = bool
    issuer  = string
  })
  default = {
    enabled = true
    issuer  = "letsencrypt-prod"
  }
}
