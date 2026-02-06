variable "name" {
  description = "Name for the Kubernetes service"
  type        = string
}

variable "namespace" {
  description = "Kubernetes namespace"
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

variable "target" {
  description = "Target deployment/pod label to route to"
  type        = string
}

variable "target_type" {
  description = "Type of target (deployment or function)"
  type        = string
  default     = "deployment"
}

variable "port" {
  description = "Service port"
  type        = number
}

variable "target_port" {
  description = "Target port on the container"
  type        = number
  default     = null
}
