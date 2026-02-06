variable "namespace" {
  description = "Kubernetes namespace"
  type        = string
}

variable "kubeconfig" {
  description = "Kubeconfig for cluster access"
  type        = string
  sensitive   = true
}

variable "name" {
  description = "Service name"
  type        = string
  default     = "service"
}

variable "target" {
  description = "Target deployment name (label selector)"
  type        = string
}

variable "target_type" {
  description = "Target type"
  type        = string
  default     = "deployment"
}

variable "port" {
  description = "Service port"
  type        = number
  default     = 80
}

variable "target_port" {
  description = "Target container port"
  type        = number
  default     = null
}
