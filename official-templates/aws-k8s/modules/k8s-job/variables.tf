variable "name" {
  description = "Job name"
  type        = string
}

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

variable "backoff_limit" {
  description = "Number of retries before marking job as failed"
  type        = number
  default     = 3
}

variable "ttl_seconds_after_finished" {
  description = "TTL for cleaning up finished jobs"
  type        = number
  default     = 300
}
