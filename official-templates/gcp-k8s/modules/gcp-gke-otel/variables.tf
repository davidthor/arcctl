variable "name" {
  description = "Name for the OTel collector deployment"
  type        = string
}

variable "namespace" {
  description = "Kubernetes namespace"
  type        = string
}

variable "project" {
  description = "GCP project ID"
  type        = string
}

variable "region" {
  description = "GCP region"
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
