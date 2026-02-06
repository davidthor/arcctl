variable "name" {
  description = "Name for the Kubernetes job"
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

variable "image" {
  description = "Container image to run"
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
  description = "Number of retries before considering the job failed"
  type        = number
  default     = 3
}

variable "ttl_seconds_after_finished" {
  description = "TTL for cleanup after job completion"
  type        = number
  default     = 300
}
