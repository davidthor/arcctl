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
  description = "Function name"
  type        = string
  default     = "function"
}

variable "environment" {
  description = "Environment variables"
  type        = map(string)
  default     = {}
}

variable "port" {
  description = "Container port"
  type        = number
  default     = 8080
}

variable "timeout" {
  description = "Request timeout in seconds"
  type        = number
  default     = 300
}

variable "memory" {
  description = "Memory allocation"
  type        = string
  default     = "256Mi"
}

variable "command" {
  description = "Container command"
  type        = list(string)
  default     = null
}

variable "cpu" {
  description = "CPU allocation"
  type        = any
  default     = null
}

variable "replicas" {
  description = "Replicas (Knative manages scaling)"
  type        = any
  default     = null
}

variable "runtime" {
  description = "Runtime (unused)"
  type        = any
  default     = null
}
