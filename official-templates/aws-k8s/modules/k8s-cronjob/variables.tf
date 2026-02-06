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
  description = "CronJob name"
  type        = string
  default     = "cronjob"
}

variable "command" {
  description = "Command to execute"
  type        = list(string)
  default     = null
}

variable "environment" {
  description = "Environment variables"
  type        = map(string)
  default     = {}
}

variable "schedule" {
  description = "Cron schedule expression"
  type        = string
}

variable "cpu" {
  description = "CPU request"
  type        = any
  default     = null
}

variable "memory" {
  description = "Memory request"
  type        = any
  default     = null
}

variable "port" {
  description = "Port (unused for cron)"
  type        = any
  default     = null
}

variable "replicas" {
  description = "Replicas (unused for cron)"
  type        = any
  default     = null
}

variable "runtime" {
  description = "Runtime (unused for cron)"
  type        = any
  default     = null
}

variable "timeout" {
  description = "Timeout"
  type        = any
  default     = null
}
