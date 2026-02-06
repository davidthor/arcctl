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

variable "image" {
  description = "Container image"
  type        = string
}

variable "command" {
  description = "Container command override"
  type        = list(string)
  default     = null
}

variable "environment" {
  description = "Environment variables"
  type        = map(string)
  default     = null
}

variable "cpu" {
  description = "CPU allocation"
  type        = string
  default     = null
}

variable "memory" {
  description = "Memory allocation"
  type        = string
  default     = null
}

variable "port" {
  description = "Container port"
  type        = number
  default     = null
}

variable "timeout" {
  description = "Request timeout in seconds"
  type        = number
  default     = 300
}

variable "concurrency" {
  description = "Max concurrent requests per container"
  type        = number
  default     = 0
}

variable "min_scale" {
  description = "Minimum number of replicas"
  type        = number
  default     = 0
}

variable "max_scale" {
  description = "Maximum number of replicas"
  type        = number
  default     = 10
}

variable "runtime" {
  description = "Runtime configuration (unused)"
  type        = any
  default     = null
}

variable "replicas" {
  description = "Replicas (unused, uses autoscaling)"
  type        = number
  default     = null
}

variable "volumes" {
  description = "Volume configuration (unused)"
  type        = any
  default     = null
}
