variable "name" {
  description = "Service name"
  type        = string
  default     = null
}

variable "namespace" {
  description = "Kubernetes namespace"
  type        = string
}

variable "kubeconfig" {
  description = "Kubernetes cluster kubeconfig content"
  type        = string
  sensitive   = true
}

variable "deployment" {
  description = "Target deployment name"
  type        = string
}

variable "port" {
  description = "Service port"
  type        = number
  default     = null
}

variable "function" {
  description = "Target function name (alternative to deployment)"
  type        = string
  default     = null
}
