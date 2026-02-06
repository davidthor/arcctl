variable "name" {
  description = "Namespace name"
  type        = string
}

variable "kubeconfig" {
  description = "Kubernetes cluster kubeconfig content"
  type        = string
  sensitive   = true
}
