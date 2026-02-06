variable "name" {
  description = "CronJob name"
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

variable "schedule" {
  description = "Cron schedule expression"
  type        = string
}

variable "image" {
  description = "Container image"
  type        = string
}

variable "command" {
  description = "Command to execute"
  type        = list(string)
  default     = null
}

variable "environment" {
  description = "Environment variables"
  type        = map(string)
  default     = null
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

variable "runtime" {
  description = "Runtime configuration (unused)"
  type        = any
  default     = null
}

variable "replicas" {
  description = "Replicas (unused)"
  type        = number
  default     = null
}

variable "port" {
  description = "Port (unused)"
  type        = number
  default     = null
}

variable "volumes" {
  description = "Volume configuration (unused)"
  type        = any
  default     = null
}
