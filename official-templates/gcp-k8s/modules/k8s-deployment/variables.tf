variable "name" {
  description = "Name for the deployment"
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
  description = "Container image to deploy"
  type        = string
}

variable "port" {
  description = "Container port"
  type        = number
  default     = null
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
  description = "CPU resource request/limit"
  type        = string
  default     = null
}

variable "memory" {
  description = "Memory resource request/limit"
  type        = string
  default     = null
}

variable "replicas" {
  description = "Number of replicas"
  type        = number
  default     = null
}
