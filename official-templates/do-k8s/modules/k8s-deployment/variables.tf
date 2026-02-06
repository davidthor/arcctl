variable "name" {
  description = "Deployment name"
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

variable "replicas" {
  description = "Number of replicas"
  type        = number
  default     = 1
}

variable "cpu" {
  description = "CPU request/limit"
  type        = string
  default     = null
}

variable "memory" {
  description = "Memory request/limit"
  type        = string
  default     = null
}

variable "port" {
  description = "Container port"
  type        = number
  default     = null
}

variable "health_check_path" {
  description = "HTTP health check path"
  type        = string
  default     = "/healthz"
}

variable "runtime" {
  description = "Runtime configuration (unused for container deployments)"
  type        = any
  default     = null
}

variable "volumes" {
  description = "Volume mounts"
  type        = any
  default     = null
}
