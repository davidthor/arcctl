variable "name" {
  description = "Secret name"
  type        = string
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

variable "data" {
  description = "Secret data as key-value pairs"
  type        = map(string)
  sensitive   = true
}
