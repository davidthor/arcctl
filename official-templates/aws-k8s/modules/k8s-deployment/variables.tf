variable "namespace" {
  description = "Kubernetes namespace"
  type        = string
}

variable "kubeconfig" {
  description = "Kubeconfig for cluster access"
  type        = string
  sensitive   = true
}

variable "image" {
  description = "Container image URI"
  type        = string
}

variable "name" {
  description = "Deployment name"
  type        = string
  default     = "deployment"
}

variable "command" {
  description = "Container command override"
  type        = list(string)
  default     = null
}

variable "environment" {
  description = "Environment variables"
  type        = map(string)
  default     = {}
}

variable "cpu" {
  description = "CPU request"
  type        = string
  default     = "250m"
}

variable "memory" {
  description = "Memory request"
  type        = string
  default     = "256Mi"
}

variable "replicas" {
  description = "Number of replicas"
  type        = number
  default     = 1
}

variable "port" {
  description = "Container port"
  type        = number
  default     = 8080
}

variable "runtime" {
  description = "Runtime (unused in container mode)"
  type        = any
  default     = null
}
